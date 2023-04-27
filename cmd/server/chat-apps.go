package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type server struct {
	pb.UnimplementedChatServiceServer
	userConnMap map[string]pb.ChatService_SendMessageServer
}

func (s *server) SendMessage(stream pb.ChatService_SendMessageServer) error {
	// Read the first message from the stream which should contain the sender's information.
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	// Add the sender's connection to the user-to-connection mapping.
	s.userConnMap[msg.GetUserFrom()] = stream

	msg.Success = true
	if err := stream.Send(msg); err != nil {
		return err
	}
	fmt.Printf("Connection added for: %s\n", msg.GetUserFrom())

	// Receive messages from the client stream
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// The client has closed the stream. Remove the sender's connection from the user-to-connection mapping.
			delete(s.userConnMap, msg.UserFrom)
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received message: %v\n", msg)

		// Get the recipient's connection from the user-to-connection mapping.
		conn, ok := s.userConnMap[msg.GetUserTo()]
		if !ok {
			// If the recipient is not connected, send an error back to the sender.
			msg.Txt = msg.GetUserTo() + " user is not connected."
			if err := stream.Send(msg); err != nil {
				return err
			}
			continue
		}

		// Forward the message to the recipient.
		if err := conn.Send(msg); err != nil {
			return err
		}

		// Wait for an ack from the recipient.
		ack, err := conn.Recv()
		if err != nil {
			return err
		}

		msg.Success = true
		// Send the ack back to the sender.
		if err := stream.Send(ack); err != nil {
			return err
		}
	}
}

func main() {
	// create gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,
			MaxConnectionAge:      30 * time.Minute,
			MaxConnectionAgeGrace: 5 * time.Minute,
			Time:                  1 * time.Minute,
			Timeout:               20 * time.Second,
		}),
	)

	// register ChatService server to gRPC server
	pb.RegisterChatServiceServer(grpcServer, &server{
		userConnMap: make(map[string]pb.ChatService_SendMessageServer),
	})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	/*

		// Create a channel to receive OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		// Start gRPC server in a goroutine
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("Failed to start gRPC server: %v", err)
			}
		}()

		// Wait for SIGINT signal
		<-sigChan

		// Stop gRPC server gracefully
		grpcServer.GracefulStop()
	*/
}

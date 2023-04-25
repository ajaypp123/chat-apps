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

	// Receive messages from the client stream
	for {
		sMsg, err := stream.Recv()
		if err == io.EOF {
			// Client stream is closed, so we're done
			return nil
		}
		if err != nil {
			return err
		}

		// update/override fromUserConn
		s.userConnMap[sMsg.GetUserFrom()] = stream
		if sMsg.GetUserTo() == "" {
			return nil
		}

		// Get the userTo's connection from the map
		conn, ok := s.userConnMap[sMsg.GetUserTo()]
		if !ok {
			return fmt.Errorf("user %s is not connected", sMsg.GetUserTo())
		}

		// Send the message to the userTo
		if err := conn.Send(sMsg); err != nil {
			return err
		}

		// Wait for ack from the userTo
		rMsg, err := conn.Recv()
		if err != nil {
			return err
		}

		// Send the ack to the fromUser
		if err := stream.Send(rMsg); err != nil {
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
	pb.RegisterChatServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

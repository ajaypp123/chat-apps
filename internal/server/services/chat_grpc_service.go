package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/kvstore"
	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type ChatServices interface {
	StartChatService(sigChan chan os.Signal)
}

type ChatServicesImpl struct {
	pb.UnimplementedChatServiceServer
}

func NewChatServices() ChatServices {
	return &ChatServicesImpl{}
}

func (chat *ChatServicesImpl) StartChatService(sigChan chan os.Signal) {
	// create gRPC server
	port := common.ConfigService().GetValue("grpc_port")
	lis, err := net.Listen("tcp", ":"+port)
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
	pb.RegisterChatServiceServer(grpcServer, &ChatServicesImpl{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
}

func (chat *ChatServicesImpl) SendMessage(stream pb.ChatService_SendMessageServer) error {
	// Read the first message from the stream which should contain the sender's information.
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	// Add the sender's connection to the user-to-connection mapping.
	chat.storeGrpcConnection(msg.GetUserFrom(), stream)

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
			chat.removeGrpcConnection(msg.UserFrom)
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received message: %v\n", msg)

		// Get the recipient's connection from the user-to-connection mapping.
		conn, err := chat.getGrpcConnection(msg.GetUserTo())
		if err != nil {
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

func (chat *ChatServicesImpl) storeGrpcConnection(user string, stream pb.ChatService_SendMessageServer) {
	jsonBytes, err := json.Marshal(stream)
	if err != nil {
		return
	}
	kvstore.Put(user, string(jsonBytes))
}

func (chat *ChatServicesImpl) removeGrpcConnection(user string) {
	kvstore.Delete(user)
}

func (chat *ChatServicesImpl) getGrpcConnection(user string) (pb.ChatService_SendMessageServer, error) {
	connStr, err := kvstore.Get(user)
	if err != nil {
		return nil, err
	}
	var conn pb.ChatService_SendMessageServer
	if err := json.Unmarshal([]byte(connStr), &conn); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return conn, nil
}

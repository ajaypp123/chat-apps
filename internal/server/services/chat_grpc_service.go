package services

import (
	"errors"
	"fmt"
	"io"

	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/common/logger"
	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
	"google.golang.org/grpc"
)

type ChatServicesImpl struct {
	pb.UnimplementedChatServiceServer
	ctx *appcontext.AppContext
}

// TODO replace with redis
var uMap map[string]pb.ChatService_SendMessageServer

func RegisterChatServices(grpcCtx *appcontext.AppContext, grpcServer *grpc.Server) {
	uMap = make(map[string]pb.ChatService_SendMessageServer)
	pb.RegisterChatServiceServer(grpcServer, &ChatServicesImpl{ctx: grpcCtx})
}

func (chat *ChatServicesImpl) SendMessage(stream pb.ChatService_SendMessageServer) error {

	if err := chat.firstConnRequest(chat.ctx, stream); err != nil {
		return err
	}

	// Receive messages from the client stream
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// The client has closed the stream. Remove the sender's connection from the user-to-connection mapping.
			chat.removeGrpcConnection(chat.ctx, msg.UserFrom)
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received message: %v\n", msg)

		// Get the recipient's connection from the user-to-connection mapping.
		conn, err := chat.getGrpcConnection(chat.ctx, msg.GetUserTo())
		if err != nil {
			// If the recipient is not connected, send an error back to the sender.
			msg.Txt = msg.GetUserTo() + " user is not connected."
			fmt.Printf("user is not connected: %v\n", msg)
			if err := stream.Send(msg); err != nil {
				return err
			}
			continue
		}
		fmt.Println("Conn Found....")

		// Forward the message to the recipient.
		if err := conn.Send(msg); err != nil {
			return err
		}
		fmt.Println("Message send to client")
	}
}

func (chat *ChatServicesImpl) firstConnRequest(ctx *appcontext.AppContext, stream pb.ChatService_SendMessageServer) error {
	// TODO verify user for first time connection

	// Read the first message from the stream which should contain the sender's information.
	msg, err := stream.Recv()
	if err != nil {
		logger.Error(ctx, "%v", err)
		return err
	}

	// Add the sender's connection to the user-to-connection mapping.
	chat.storeGrpcConnection(ctx, msg.GetUserFrom(), stream)

	msg.Success = true
	if err := stream.Send(msg); err != nil {
		return err
	}
	logger.Info(ctx, "Connection added for: %s\n", msg.GetUserFrom())
	return nil
}

func (chat *ChatServicesImpl) storeGrpcConnection(_ *appcontext.AppContext, user string, stream pb.ChatService_SendMessageServer) {
	uMap[user] = stream
}

func (chat *ChatServicesImpl) removeGrpcConnection(_ *appcontext.AppContext, user string) {
	delete(uMap, user)
}

func (chat *ChatServicesImpl) getGrpcConnection(_ *appcontext.AppContext, user string) (pb.ChatService_SendMessageServer, error) {
	//return conn, nil
	stream, ok := uMap[user]
	if !ok {
		return nil, errors.New("not found")
	}
	return stream, nil
}

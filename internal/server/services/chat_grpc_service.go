package services

import (
	"encoding/json"
	"errors"
	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/common/kvstore"
	"github.com/ajaypp123/chat-apps/common/logger"
	"github.com/ajaypp123/chat-apps/common/streamer"
	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
	"github.com/ajaypp123/chat-apps/internal/server/models"
	"google.golang.org/grpc"
	"io"
)

type ChatServicesImpl struct {
	pb.UnimplementedChatServiceServer
	ctx    *appcontext.AppContext
	kvs    kvstore.KVStoreI
	prefix string
	stream streamer.StreamingServiceI
}

type ChatSub struct {
	ctx *appcontext.AppContext
}

var uMap map[string]pb.ChatService_SendMessageServer

func GetSubscriber(ctx *appcontext.AppContext) streamer.Subscriber {
	return &ChatSub{ctx: ctx}
}

func RegisterChatServices(grpcCtx *appcontext.AppContext, grpcServer *grpc.Server) error {
	kv, err := kvstore.GetKVStore()
	if err != nil {
		return err
	}

	stream, err := streamer.GetStreamingService()
	if err != nil {
		return err
	}

	uMap = make(map[string]pb.ChatService_SendMessageServer)
	pb.RegisterChatServiceServer(grpcServer, &ChatServicesImpl{
		ctx:    grpcCtx,
		kvs:    kv,
		prefix: "chat-apps/cache/connection/",
		stream: stream,
	})

	topic := common.ConfigService().GetValue("streaming.topic")
	partition := common.ConfigService().GetValue("streaming.partition")
	sub := GetSubscriber(grpcCtx)
	stream.RegisterSubscriber(sub)
	err = stream.StartListening(topic, partition)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatSub) ReceiveMsg(msg interface{}) {
	ctx := c.ctx
	logMsg := "ChatSub::ReceiveMsg "
	logger.Info(ctx, logMsg, "entry")

	var pbMsg pb.Meg
	err := json.Unmarshal([]byte(msg.(string)), &pbMsg)
	if err != nil {
		logger.Error(ctx, logMsg, "Failed to convert message ", msg, err)
		return
	}

	stream, ok := uMap[pbMsg.GetUserTo()]
	if !ok {
		logger.Error(ctx, logMsg, "user: ", pbMsg.GetUserTo(), " may be not connected, add message in database...")
		// TODO: handle close connections
		return
	}
	err = stream.Send(&pbMsg)
	if err != nil {
		logger.Error(ctx, logMsg, "failed to send message, err: ", err)
		return
	}
}

func (chat *ChatServicesImpl) SendMessage(stream pb.ChatService_SendMessageServer) error {
	ctx := chat.ctx
	logMsg := "ChatServicesImpl::SendMessage "
	logger.Debug(ctx, "entry", logMsg)
	errCount := 0

	done := stream.Context().Done()

	//go func() {
	//	<-done
	//	logger.Info(ctx, "StreamMessages stream closed: ", ctx.Err())
	//	// Close the message channel to signal the loop to exit
	//}()

	// Receive messages from the client stream
	for {
		if errCount == 5 {
			break
		}

		select {
		case <-done:
			logger.Info(ctx, logMsg, "Stream is closed")
			return nil
		default:
			msg, err := chat.receive(ctx, stream)
			if err == io.EOF {
				// stream is closed
				return nil
			}
			if err != nil {
				logger.Warn(ctx, logMsg, " Failed to parse message: ", msg, err)
				continue
			}

			// publish message
			logger.Info(ctx, logMsg, " publish message ", msg)
			ok, err := chat.publishMessage(ctx, msg)
			if err != nil {
				errCount = errCount + 1
			}
			if !ok {
				logger.Warn(ctx, logMsg, " Failed to publish message: ", msg)
			}
		}
	}

	return nil
}

func (chat *ChatServicesImpl) receive(ctx *appcontext.AppContext, stream pb.ChatService_SendMessageServer) (*pb.Meg, error) {
	logMsg := "ChatServicesImpl::receive "
	logger.Debug(ctx, logMsg, "entry")

	msg, err := stream.Recv()
	if err == io.EOF {
		// The client has closed the stream. Remove the sender's connection from the user-to-connection mapping.
		logger.Info(ctx, logMsg, "Closing connection stream from "+msg.GetUserFrom())
		if err := chat.removeGrpcConnection(chat.ctx, msg.UserFrom); err != nil {
			return msg, err
		}
		return msg, err
	}
	if err != nil {
		logger.Error(ctx, logMsg, "Err while receiving message ", err)
		return msg, err
	}

	err = chat.storeGrpcConnection(ctx, msg.UserFrom, stream)
	if err != nil {
		return msg, err
	}
	logger.Debug(ctx, logMsg, "Received message: ", msg)
	return msg, err
}

func (chat *ChatServicesImpl) storeGrpcConnection(ctx *appcontext.AppContext, user string, stream pb.ChatService_SendMessageServer) error {
	if _, ok := uMap[user]; !ok {
		conn := &models.ConnectionDetail{
			Topic:     common.ConfigService().GetValue("streaming.topic"),
			Partition: common.ConfigService().GetValue("streaming.partition"),
			Username:  user,
		}
		data, err := json.Marshal(conn)
		if err != nil {
			logger.Error(ctx, "ChatServicesImpl::storeGrpcConnection failed to store key, err:", err)
			return err
		}
		err = chat.kvs.Put(chat.prefix+user, string(data))
		if err != nil {
			return err
		}
	}
	uMap[user] = stream
	return nil
}

func (chat *ChatServicesImpl) removeGrpcConnection(_ *appcontext.AppContext, user string) error {
	delete(uMap, user)
	return chat.kvs.Delete(chat.prefix + user)
}

func (chat *ChatServicesImpl) getGrpcConnection(_ *appcontext.AppContext, user string) (pb.ChatService_SendMessageServer, error) {
	stream, ok := uMap[user]
	if !ok {
		return nil, errors.New("not found")
	}
	return stream, nil
}

func (chat *ChatServicesImpl) publishMessage(ctx *appcontext.AppContext, msg *pb.Meg) (bool, error) {
	logMsg := "ChatServicesImpl:: publishMessage "

	// publish message
	strConn, err := chat.kvs.Get(chat.prefix + msg.GetUserTo())
	if err != nil {
		logger.Error(ctx, logMsg, "failed to publish message. Unknown sender, err: ", err, msg)
		return false, nil
	}
	var conn models.ConnectionDetail
	if err := json.Unmarshal([]byte(strConn), &conn); err != nil {
		logger.Error(ctx, "failed to parse message. err: ", err)
		return false, nil
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		logger.Error(ctx, logMsg, "failed to parse message, err: ", err, msg)
		return false, nil
	}

	err = chat.stream.PublishMessage(conn.Topic, conn.Partition, string(msgByte))
	if err != nil {
		logger.Error(ctx, logMsg, "failed to publish message, err: ", err, msg)
		return false, err
	}
	return true, nil
}

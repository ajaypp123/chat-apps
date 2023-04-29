package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/kvstore"
	"github.com/ajaypp123/chat-apps/common/logger"
	"github.com/ajaypp123/chat-apps/internal/server/services"
)

func main() {
	// Create a channel to receive OS signals
	fmt.Println("Register signal with service")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// grpc service
	var chat services.ChatServices
	{
		chat = services.NewChatServices()
	}
	chat.StartChatService(sigChan)

	// Wait for SIGINT signal
	<-sigChan
}

// init is self called and will initialise all services
func init() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "index", "server")

	// log initalized
	err := logger.NewLogger(ctx, "server.log", logger.DEBUG)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	// logger.Info(ctx, "Starting application...")

	// config
	if err := common.ConfigService().Init(); err != nil {
		logger.Error(ctx, "Failed to setup config, exit from service. err: %v", err)
		os.Exit(1)
	}

	// fmt.Println(common.ConfigService().GetValue("port"))

	// kvstore
	if err := kvstore.Init("mem"); err != nil {
		logger.Error(ctx, "Failed to setup kvstore. err: ", err)
		os.Exit(1)
	}
}

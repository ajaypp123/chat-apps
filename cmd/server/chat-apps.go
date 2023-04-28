package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajaypp123/chat-apps/common"
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
	if err := common.ConfigService().Init(); err != nil {
		log.Fatalf("Failed to setup config, exit from service. err: %v", err)
		os.Exit(1)
	}
	fmt.Println(common.ConfigService().GetValue("port"))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/common/kvstore"
	"github.com/ajaypp123/chat-apps/common/logger"
	"github.com/ajaypp123/chat-apps/internal/server/controller"
	"github.com/ajaypp123/chat-apps/internal/server/services"
	"github.com/gorilla/mux"
)

func main() {
	ctx := appcontext.DefaultContext()
	logger.Info(ctx, "Server started at : ", common.GetTimeNow())

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// grpc service
	var chatService services.ChatServices
	{
		chatService = services.NewChatServices()
		chatService.StartChatService(sigChan)
		defer chatService.Close()
	}

	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// define REST endpoints
	controller.NewUserController(services.NewUserService()).RegisterUserHandler(router)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
	// Wait for SIGINT signal
	<-sigChan
}

// init is self called and will initialise all services
func init() {
	ctx := appcontext.DefaultContext()

	// log initalized
	err := logger.NewLogger(ctx, "server.log", logger.DEBUG)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	if err := common.ConfigService().Init(); err != nil {
		logger.Error(ctx, "Failed to setup config, exit from service. err: %v", err)
		os.Exit(1)
	}

	// kvstore
	if err := kvstore.Init("mem"); err != nil {
		logger.Error(ctx, "Failed to setup kvstore. err: ", err)
		os.Exit(1)
	}
}

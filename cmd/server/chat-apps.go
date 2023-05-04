package main

import (
	"flag"
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

var (
	httpPort *string
	grpcPort *string
)

func main() {
	httpPort = flag.String("http", ":8080", "HTTP server address")
	grpcPort = flag.String("grpc", ":50051", "Grpc server address")
	flag.Parse()

	ctx := appcontext.GetDefaultContext()
	ctx.AddValue("index", "server")
	err := logger.NewLogger(ctx, "server.log", logger.DEBUG)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	logger.Info(ctx, "Server started at :", common.GetTimeNow())

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	// grpc service
	grpc_service := common.GrpcHelper{}
	grpcServer := grpc_service.SetGrpcServer(*grpcPort)
	services.RegisterChatServices(grpcServer)
	grpc_service.StartGrpcServer()

	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    *httpPort,
		Handler: router,
	}

	// define REST endpoints
	controller.NewUserController(services.NewUserService()).RegisterUserHandler(router)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}

	// Wait for SIGINT signal
	<-sigChan
	fmt.Println("Exit from:", *httpPort, *grpcPort)
	grpc_service.StopGrpcServer()
	logger.Close(ctx)
}

// init is self called and will initialise all services
func init() {
	if err := common.ConfigService().Init(); err != nil {
		panic(fmt.Sprintf("Failed to setup config, exit from service. err: %v", err))
	}

	// kvstore
	if err := kvstore.Init("mem"); err != nil {
		panic(fmt.Sprintf("Failed to setup kvstore. err: %v", err))
	}
}

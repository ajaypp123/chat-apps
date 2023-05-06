package main

import (
	"flag"
	"fmt"
	"github.com/ajaypp123/chat-apps/common/streamer"
	"github.com/ajaypp123/chat-apps/internal/server/repos"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/appcontext"
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

	kv := make(map[string]string)
	kv["index"] = "server"
	kv["process"] = "chat-apps"
	kv["id"] = common.ConfigService().GetValue("server_id")
	ctx := appcontext.GetDefaultContext(kv)
	err := logger.NewLogger(ctx, "server.log", logger.DEBUG)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	logger.Info(ctx, "Server started at :", common.GetTimeNow())

	// grpc service
	grpcCtx := ctx.DeepCopy()
	grpcCtx.AddValue("process", "grpc")

	grpcService := common.GrpcHelper{}
	grpcServer := grpcService.SetGrpcServer(ctx, *grpcPort)
	err = services.RegisterChatServices(grpcCtx, grpcServer)
	if err != nil {
		panic(err)
	}
	grpcService.StartGrpcServer(ctx)

	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    *httpPort,
		Handler: router,
	}

	// define REST endpoints
	controller.NewHealthController().RegisterHealthHandler(router)
	controller.NewUserController(services.NewUserService()).RegisterUserHandler(router)

	// Instantiate server, and listen on specified port.
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- httpServer.ListenAndServe()
	}()

	// Wait for SIGINT signal
	logger.Debug(ctx, "Exit from server", <-errs)
	grpcService.StopGrpcServer(ctx)
	steam, _ := streamer.GetStreamingService()
	steam.StopListening()
	logger.Info(ctx, "Exit from application ....")
	logger.Close(ctx)
}

// init is self called and will initialise all services
func init() {
	if err := common.ConfigService().Init(); err != nil {
		panic(fmt.Sprintf("Failed to setup config, exit from service. err: %v", err))
	}

	if err := repos.InitializeDB(repos.RedisDB); err != nil {
		panic(fmt.Sprintf("Failed to setup kvstore. err: %v", err))
	}

	addr := common.ConfigService().GetValue("redis.addr")
	pass := common.ConfigService().GetValue("redis.pass")
	db := 1 // 0 default and 1 for connection
	if _, err := streamer.NewRedisStreamingService(addr, pass, db); err != nil {
		panic(fmt.Sprintf("Failed to setup streamer. err: %v", err))
	}
}

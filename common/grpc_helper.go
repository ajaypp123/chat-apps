package common

import (
	"fmt"
	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/common/logger"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	listener   net.Listener = nil
	GrpcServer *grpc.Server = nil
)

type GrpcHelper struct{}

func (g *GrpcHelper) SetGrpcServer(ctx *appcontext.AppContext, port string) *grpc.Server {
	if listener == nil {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			logger.Error(ctx, "failed to listen to port ", port, err)
			os.Exit(1)
		}
		listener = lis
	}

	if GrpcServer == nil {
		GrpcServer = grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     5 * time.Minute,
				MaxConnectionAge:      30 * time.Minute,
				MaxConnectionAgeGrace: 5 * time.Minute,
				Time:                  1 * time.Minute,
				Timeout:               3 * time.Second,
			}),
		)
	}
	return GrpcServer
}

func (g *GrpcHelper) StartGrpcServer(ctx *appcontext.AppContext) {
	go func() {
		if err := GrpcServer.Serve(listener); err != nil {
			logger.Error(ctx, "Failed to start gRPC server: ", err)
		}
		fmt.Println("stop")
	}()
}

func (g *GrpcHelper) StopGrpcServer(ctx *appcontext.AppContext) {
	logger.Info(ctx, "Stopping grpc server")
	GrpcServer.GracefulStop()
}

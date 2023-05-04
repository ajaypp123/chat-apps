package common

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	listener   net.Listener = nil
	GrpcServer *grpc.Server = nil
)

type GrpcHelper struct{}

func (g *GrpcHelper) SetGrpcServer(port string) *grpc.Server {
	if listener == nil {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
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

func (g *GrpcHelper) StartGrpcServer() {
	go func() {
		if err := GrpcServer.Serve(listener); err != nil {
			panic(fmt.Sprintf("Failed to start gRPC server: %v", err))
		}
	}()
}

func (g *GrpcHelper) StopGrpcServer() {
	GrpcServer.GracefulStop()
}

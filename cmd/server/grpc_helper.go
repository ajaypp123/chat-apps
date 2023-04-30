package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func setGrpcServer(port string) {
	if listener == nil {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		listener = lis
	}

	if grpcServer == nil {
		grpcServer = grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle:     5 * time.Minute,
				MaxConnectionAge:      30 * time.Minute,
				MaxConnectionAgeGrace: 5 * time.Minute,
				Time:                  1 * time.Minute,
				Timeout:               3 * time.Second,
			}),
		)
	}
}

func startGrpcServer() {
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
	defer grpcServer.GracefulStop()
}

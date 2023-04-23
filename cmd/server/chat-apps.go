package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) SendMessage(ctx context.Context, msg *pb.Meg) (*pb.Ack, error) {
	// add your SendMessage logic here
	return nil, nil
}

func main() {
	// web connection
	// create gRPC server
	grpcServer := grpc.NewServer()

	// register ChatService server to gRPC server
	pb.RegisterChatServiceServer(grpcServer, &server{})

	// create Gorilla Mux router
	router := mux.NewRouter()

	// attach gRPC server to Gorilla Mux server
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		grpcServer.ServeHTTP(w, r)
	})

	// start Gorilla Mux server
	log.Fatal(http.ListenAndServe(":8080", router))
}

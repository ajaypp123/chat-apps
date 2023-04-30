package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ajaypp123/chat-apps/internal/client/services"
	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Resp struct {
	ReqID  string `json:"req_id"`
	Status string `json:"status"`
	Data   struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Secret   string `json:"secret"`
	} `json:"data"`
	Code int `json:"code"`
}

func main() {

	grpcPort := flag.String("grpc", ":50051", "Grpc server address")
	httpPort := flag.String("http", ":8080", "HTTP server address")
	flag.Parse()
	conn, err := grpc.Dial("localhost"+*grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to connect to server: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	stream, err := client.SendMessage(context.Background())
	if err != nil {
		fmt.Println("Error sending message: ", err)
		os.Exit(1)
	}

	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	fmt.Println("Fetching Data")
	services.FetchUserData(*httpPort, user)

	connectMsg := &pb.Meg{
		Id:       "0",
		UserFrom: user,
		UserTo:   "",
		Txt:      "connect",
		Success:  false,
	}

	// Send the message to the server and wait for a response.
	if err := stream.Send(connectMsg); err != nil {
		fmt.Println("Error sending message: ", err)
		os.Exit(1)
	}

	// Wait for an ack from the server
	ack, err := stream.Recv()
	if err != nil {
		fmt.Println("failed to receive ack: ", err)
		os.Exit(1)
	}

	fmt.Printf("Connection Request: User %s, registered successfuly with server. msg: %v \n\n\n", user, ack)

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				fmt.Println("Error receiving message: ", err)
			}
			if msg.UserTo == user {
				services.AddNewMessage(*httpPort, msg.UserFrom, msg.Txt)
				//fmt.Printf("[%s]: %s\n", msg.UserFrom, msg.Txt)
			}
		}
	}()

	for i := 0; i < 100; i++ {
		fmt.Printf("\n\n\nUsername: %s\n\n", services.Data.Username)
		services.PrintTable()
		fmt.Printf("\n\n\n")

		msg := &pb.Meg{
			Id:       strconv.Itoa(i),
			UserFrom: "user1",
			UserTo:   "user1",
			Txt:      "Hiii " + strconv.Itoa(i),
			Success:  false,
		}

		fmt.Printf("\n%s:%s - %s\n", msg.UserFrom, msg.UserTo, msg.Txt)

		// Send the message to the server and wait for a response.
		if err := stream.Send(msg); err != nil {
			fmt.Println("Error sending message: ", err)
		}
		fmt.Println("Message send done....")

	}
}

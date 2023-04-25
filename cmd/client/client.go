package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	stream, err := client.SendMessage(context.Background())
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	var user string
	var count int64 = 1
	// get user name
	fmt.Println("Enter Your First Name: ")
	fmt.Scanln(&user)

	var sender string
	for {
		msgId := strconv.FormatInt(count, 10)

		// get sender name
		fmt.Println("Enter Sender Name: ")
		fmt.Scanln(&sender)

		msg := &pb.Meg{
			Id:       msgId,
			UserFrom: user,
			UserTo:   sender,
			Txt:      "Hi for " + msgId,
			Success:  false,
		}

		fmt.Printf("%s:%s - %s", msg.UserFrom, msg.UserTo, msg.Txt)

		// Send the message to the server and wait for a response.
		if err := stream.Send(msg); err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		// Receive and print the server's response messages until the stream is closed.
		for {
			resp, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error receiving message: %v", err)
			}

			if err == io.EOF {
				break
			}

			log.Printf("Received message: %v", resp)
			if !resp.Success {
				log.Fatalf("Error message not received: %v", err)
			}
			fmt.Printf("%s:%s - %s", resp.UserFrom, resp.UserTo, resp.Txt)
		}
		count += 1
	}
}

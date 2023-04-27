package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	var count int64 = 1
	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	connectMsg := &pb.Meg{
		Id:       "0",
		UserFrom: user,
		UserTo:   "",
		Txt:      "connect",
		Success:  false,
	}

	// Send the message to the server and wait for a response.
	if err := stream.Send(connectMsg); err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	// Wait for an ack from the server
	ack, err := stream.Recv()
	if err != nil {
		log.Fatalf("failed to receive ack: %v", err)
	}

	fmt.Printf("User %s, registered successfuly with server. msg: %v \n", user, ack)

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error receiving message: %v", err)
			}
			if msg.UserTo == user {
				fmt.Printf("[%s]: %s\n", msg.UserFrom, msg.Txt)
			}
		}
	}()

	for {
		msgId := strconv.FormatInt(count, 10)

		fmt.Print("Enter recipient name: ")
		toUser, _ := reader.ReadString('\n')
		toUser = strings.TrimSpace(toUser)

		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		msg := &pb.Meg{
			Id:       msgId,
			UserFrom: user,
			UserTo:   toUser,
			Txt:      message,
			Success:  false,
		}

		fmt.Printf("%s:%s - %s\n", msg.UserFrom, msg.UserTo, msg.Txt)

		// Send the message to the server and wait for a response.
		if err := stream.Send(msg); err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		// Wait for an ack from the server
		ack, err := stream.Recv()
		if err != nil {
			log.Fatalf("failed to receive ack: %v", err)
		}

		if !ack.Success {
			fmt.Printf("Message sent failed. ID: %s, %v\n", msgId, ack)
		} else {
			fmt.Printf("Message sent successfully. ID: %s, %v\n", msgId, ack)
		}
		count += 1
	}
}

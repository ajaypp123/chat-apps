package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
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

	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	fmt.Println("Fetching Data")
	services.FetchUserData(*httpPort, user)

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

	go func() {
		<-stream.Context().Done()
		fmt.Println("Stream closed")
	}()

	var count int64 = 1

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

	/*// Wait for an ack from the server
	ack, err := stream.Recv()
	if err != nil {
		fmt.Println("failed to receive ack: ", err)
		os.Exit(1)
	}*/

	fmt.Printf("Connection Request: User %s, registered successfuly with server \n\n\n", connectMsg)

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				md, er := stream.Header()
				fmt.Println("Error receiving message: ", er, err)
				for k, v := range md {
					log.Printf("Metadata[%s]: %s", k, v)
				}
				fmt.Println("Error receiving message: ", err)
			}
			if msg.UserTo == user {
				services.AddNewMessage(*httpPort, msg.UserFrom, msg.Txt)
				//fmt.Printf("[%s]: %s\n", msg.UserFrom, msg.Txt)
			}
		}
	}()

	for {
		fmt.Printf("\n\n\nUsername: %s\n\n", services.Data.Username)
		services.PrintTable()
		fmt.Printf("\n\n\n")

		fmt.Println(`Enter your choise: 
			a. Send message
			b. Show perticular user message
			c. Refresh
			d. Create Group *TBD
			e. Exit`)
		choise, _ := reader.ReadString('\n')
		choise = strings.TrimSpace(choise)

		switch choise {
		case "a":
			msgId := strconv.FormatInt(count, 10)
			services.SendMessage(msgId, reader, stream)
			count += 1
			fmt.Println("Message send completed")
		case "b":
			services.ShowUserMessage(reader)
		case "c":
			continue
		}

	}
}

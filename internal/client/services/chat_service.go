package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	pb "github.com/ajaypp123/chat-apps/internal/communication_grpc"
)

func FetchUserData(port, username string) {
	resp, err := http.Get(fmt.Sprintf("http://localhost"+port+"/v1/chatapp/users?username=%s", username))
	if err != nil {
		fmt.Printf("failed to get user info: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var record Resp
	err = json.NewDecoder(resp.Body).Decode(&record)
	if err != nil {
		fmt.Println("failed to get user info: ", err)
		os.Exit(1)
	}

	Data = ClientData{}
	Data.Name = record.Data.Name
	Data.Username = record.Data.Username
	Data.Phone = record.Data.Phone
	Data.Secret = record.Data.Secret
	fmt.Println("your identity: ", Data)
}

func ShowUserMessage(reader *bufio.Reader) {
	fmt.Print("\n\nEnter recipient username: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	ShowMessage(user)
}

func SendMessage(msgId string, reader *bufio.Reader, stream pb.ChatService_SendMessageClient) {
	fmt.Print("Enter recipient username: ")
	toUser, _ := reader.ReadString('\n')
	toUser = strings.TrimSpace(toUser)

	fmt.Print("Enter message: ")
	message, _ := reader.ReadString('\n')
	message = strings.TrimSpace(message)

	msg := &pb.Meg{
		Id:       msgId,
		UserFrom: Data.Username,
		UserTo:   toUser,
		Txt:      message,
		Success:  false,
	}

	fmt.Printf("\n%s:%s - %s\n", msg.UserFrom, msg.UserTo, msg.Txt)

	// Send the message to the server and wait for a response.
	if err := stream.Send(msg); err != nil {
		fmt.Println("Error sending message: ", err)
	}
	fmt.Println("Message send done....")
}

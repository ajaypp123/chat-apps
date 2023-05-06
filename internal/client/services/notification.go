package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type notification struct {
	username string
	phone    string
	queue    []string
}

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

var (
	table     = make(map[string]*notification)
	tableLock sync.Mutex
)

func ShowMessage(username string) {
	n, ok := table[username]
	if !ok {
		fmt.Println("No New Message")
		return
	}
	size := len(n.queue)
	for i := 0; i < size; i++ {
		fmt.Printf("%s - %s : %s\n", Data.Username, username, n.queue[0])
		n.queue = n.queue[1:]
	}
}

func AddNewMessage(port, username string, message string) error {
	tableLock.Lock()
	defer tableLock.Unlock()

	n, ok := table[username]
	if !ok {
		// User not in table, get their info from API
		resp, err := http.Get(fmt.Sprintf("http://localhost"+port+"/v1/chat-apps/users?username=%s", username))
		if err != nil {
			return fmt.Errorf("failed to get user info: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get user info, status code: %d", resp.StatusCode)
		}

		var data *Resp
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode user info: %v", err)
		}

		n = &notification{
			username: data.Data.Username,
			phone:    data.Data.Phone,
			queue:    make([]string, 0),
		}

		table[username] = n
	}

	n.queue = append(n.queue, message)
	fmt.Println("Notification added...")
	return nil
}

func PrintTable() {

	fmt.Printf("%-15s | %-15s | %-15s\n", "username", "phone", "notification")
	fmt.Println("-----------------------------------------------------")

	for _, n := range table {
		fmt.Printf("%-15v | %-15v | %-15v\n", n.username, n.phone, strconv.Itoa(len(n.queue)))
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func main() {
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Create subscriber
	pubsub := rdb.Subscribe(context.Background(), "example-channel")
	defer pubsub.Close()

	// Subscribe to channel
	ch := pubsub.Channel()

	// Publish messages to channel every 1 second
	go func() {
		for i := 1; i <= 10; i++ {
			err := rdb.Publish(context.Background(), "example-channel", fmt.Sprintf("Message %d", i)).Err()
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Consume messages from channel
	for msg := range ch {
		fmt.Println("Received message:", msg.Payload)
	}
}

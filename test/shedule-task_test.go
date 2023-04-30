package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRedisKVStore(t *testing.T) {
	// Create a ticker that ticks at a fixed interval
	ticker := time.NewTicker(2 * time.Second)

	// Continuously check the connection status and reconnect if necessary
	go func() {
		for range ticker.C {
			fmt.Println("in loop")
		}
	}()

	for {
		fmt.Println("outside")
		time.Sleep(time.Minute * 1)
	}
}

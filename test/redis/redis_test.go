package main

import (
	"testing"
	"time"

	"github.com/ajaypp123/chat-apps/common/kvstore"
)

func TestRedisKVStore(t *testing.T) {
	store, err := kvstore.NewRedisKVStore("localhost:6379", "", 0)
	if err != nil {
		t.Fatalf("failed to create RedisKVStore: %v", err)
	}
	//defer store.Close()

	// Test Put, Get, and Delete
	key := "test"
	value := "123"
	if err := store.Put(key, value); err != nil {
		t.Fatalf("failed to Put value: %v", err)
	}
	got, err := store.Get(key)
	if err != nil {
		t.Fatalf("failed to Get value: %v", err)
	}
	if got != value {
		t.Fatalf("expected %q, got %q", value, got)
	}
	if err := store.Delete(key); err != nil {
		t.Fatalf("failed to Delete value: %v", err)
	}

	// Test Lock and Unlock
	key = "lock_test"
	timeout := 5 * time.Second
	ok, err := store.Lock(key, timeout)
	if err != nil {
		t.Fatalf("failed to Lock key: %v", err)
	}
	if !ok {
		t.Fatalf("expected to Lock key")
	}
	// Test that another Lock on the same key fails
	ok, err = store.Lock(key, timeout)
	if err != nil {
		t.Fatalf("failed to Lock key: %v", err)
	}
	if ok {
		t.Fatalf("expected Lock to fail")
	}
	if err := store.Unlock(key); err != nil {
		t.Fatalf("failed to Unlock key: %v", err)
	}
}

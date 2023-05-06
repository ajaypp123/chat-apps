package kvstore

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var ctx = context.Background()

// RedisKVStore implements the KVStoreI interface
type RedisKVStore struct {
	client *redis.Client
	mu     sync.Mutex
}

// NewRedisKVStore creates a new RedisKVStore instance
func NewRedisKVStore(addr string, password string, db int) (KVStoreI, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Check if Redis is connected
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	cli := &RedisKVStore{client: client}
	if defaultClient == nil {
		defaultClient = cli
	}

	return cli, nil
}

// Get retrieves the value for a given key from Redis
func (s *RedisKVStore) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("key not found")
	} else if err != nil {
		return "", err
	}

	return value, nil
}

// Put adds or updates the value for a given key in Redis
func (s *RedisKVStore) Put(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.client.Set(ctx, key, value, 0).Err()
}

// Delete deletes the value for a given key from Redis
func (s *RedisKVStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.client.Del(ctx, key).Err()
}

// Lock obtains a lock on a given key using Redis
func (s *RedisKVStore) Lock(key string, timeout time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.client.SetNX(ctx, key, "locked", timeout).Result()
}

// Unlock releases a lock on a given key using Redis
func (s *RedisKVStore) Unlock(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.client.Del(ctx, key).Err()
}

package kvstore

import (
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var defaultClient *RedisKVStore

// RedisKVStore implements the KVStoreI interface
type RedisKVStore struct {
	client *redis.Client
	mu     sync.Mutex
}

// NewRedisKVStore creates a new RedisKVStore instance
func NewRedisKVStore(addr string, password string, db int) (KVStoreI, error) {
	if defaultClient != nil {
		return defaultClient, nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Check if Redis is connected
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	defaultClient = &RedisKVStore{client: client}
	return defaultClient, nil
}

func GetRedisKVStore() (KVStoreI, error) {
	if defaultClient != nil {
		return defaultClient, nil
	}
	return nil, errors.New("KVStore is not initialised")
}

// Get retrieves the value for a given key from Redis
func (s *RedisKVStore) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, err := s.client.Get(key).Result()
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

	return s.client.Set(key, value, 0).Err()
}

// Delete deletes the value for a given key from Redis
func (s *RedisKVStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.client.Del(key).Err()
}

// Lock obtains a lock on a given key using Redis
func (s *RedisKVStore) Lock(key string, timeout time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.client.SetNX(key, "locked", timeout).Result()
}

// Unlock releases a lock on a given key using Redis
func (s *RedisKVStore) Unlock(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.client.Del(key).Err()
}

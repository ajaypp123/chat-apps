package common

/*
import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

type RedisKVStore struct {
	client *redis.Client
}

func NewRedisKVStore() (*RedisKVStore, error) {
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Test connection to Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisKVStore{client: client}, nil
}

func (m *RedisKVStore) Get(key interface{}) (interface{}, error) {
	// Acquire lock
	lockKey := m.lockKey(key)
	lock, err := m.client.SetNX(context.Background(), lockKey, "1", 1*time.Second).Result()
	if err != nil {
		return nil, err
	}
	if !lock {
		return nil, errors.New("could not acquire lock")
	}
	defer m.client.Del(context.Background(), lockKey)

	// Get value
	value, err := m.client.Get(context.Background(), key.(string)).Result()
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (m *RedisKVStore) Put(key interface{}, value interface{}) error {
	// Acquire lock
	lockKey := m.lockKey(key)
	lock, err := m.client.SetNX(context.Background(), lockKey, "1", 1*time.Second).Result()
	if err != nil {
		return err
	}
	if !lock {
		return errors.New("could not acquire lock")
	}
	defer m.client.Del(context.Background(), lockKey)

	// Put value
	return m.client.Set(context.Background(), key.(string), value.(string), 0).Err()
}

func (m *RedisKVStore) Delete(key interface{}) error {
	// Acquire lock
	lockKey := m.lockKey(key)
	lock, err := m.client.SetNX(context.Background(), lockKey, "1", 1*time.Second).Result()
	if err != nil {
		return err
	}
	if !lock {
		return errors.New("could not acquire lock")
	}
	defer m.client.Del(context.Background(), lockKey)
	return nil
}
*/

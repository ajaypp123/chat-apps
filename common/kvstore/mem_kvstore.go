package kvstore

import (
	"errors"
	"time"
)

type MemKVStore struct {
	store map[string]string
}

func GetMemKVStore() KVStoreI {
	return &MemKVStore{
		store: make(map[string]string),
	}
}

func (s *MemKVStore) Get(key string) (string, error) {
	value, ok := s.store[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (s *MemKVStore) Put(key string, value string) error {
	s.store[key] = value
	return nil
}

func (s *MemKVStore) Delete(key string) error {
	delete(s.store, key)
	return nil
}

func (s *MemKVStore) Lock(key string, timeout time.Duration) (bool, error) {
	return false, nil
}

func (s *MemKVStore) Unlock(key string) error {
	return nil
}

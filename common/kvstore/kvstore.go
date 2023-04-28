package common

import "fmt"

type KVStoreI interface {
	Get(key interface{}) (interface{}, error)
	Put(key interface{}, value interface{}) error
	Delete(key interface{}) error
}

type KVStore struct {
	store map[interface{}]interface{}
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[interface{}]interface{}),
	}
}

func (s *KVStore) Get(key interface{}) (interface{}, error) {
	value, ok := s.store[key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	return value, nil
}

func (s *KVStore) Put(key interface{}, value interface{}) error {
	s.store[key] = value
	return nil
}

func (s *KVStore) Delete(key interface{}) error {
	delete(s.store, key)
	return nil
}

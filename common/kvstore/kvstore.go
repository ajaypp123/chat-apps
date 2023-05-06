package kvstore

import (
	"errors"
	"time"
)

type KVStoreI interface {
	Get(key string) (string, error)
	Put(key string, value string) error
	Delete(key string) error
	Lock(key string, timeout time.Duration) (bool, error)
	Unlock(key string) error
}

var defaultClient KVStoreI

func GetKVStore() (KVStoreI, error) {
	if defaultClient != nil {
		return defaultClient, nil
	}
	return nil, errors.New("KVStore is not initialised")
}

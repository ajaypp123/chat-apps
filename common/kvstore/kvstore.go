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

var kv KVStoreI = nil

// Init name for memory is mem and redis for redis
// this is used to set default kv store
func Init(name string) error {
	if name == "mem" {
		kv = GetMemKVStore()
		return nil
	}
	return errors.New("invalid kv name")
}

func Get(key string) (string, error) {
	return kv.Get(key)
}

func Put(key string, value string) error {
	return kv.Put(key, value)
}

func Delete(key string) error {
	return kv.Delete(key)
}

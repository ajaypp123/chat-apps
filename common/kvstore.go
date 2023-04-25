package common

type KVStoreI interface{}

type KVStore struct {
}

func GetKVStore() KVStoreI {
	return &KVStore{}
}

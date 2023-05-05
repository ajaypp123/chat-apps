package repos

import (
	"errors"
	"fmt"
	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/kvstore"
	"strconv"
)

const (
	RedisDB string = "redis"
)

func initializeRedis() error {
	// kvstore
	addr := common.ConfigService().GetValue("redis.addr")
	pass := common.ConfigService().GetValue("redis.pass")
	db := common.ConfigService().GetValue("redis.db")
	dbVal, err := strconv.Atoi(db)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to setup redis, exit from service. err: %v", err))
	}
	if _, err := kvstore.NewRedisKVStore(addr, pass, dbVal); err != nil {
		return errors.New(fmt.Sprintf("Failed to setup kvstore. err: %v", err))
	}
	return nil
}

func InitializeDB(DBName string) error {
	switch DBName {
	case RedisDB:
		return initializeRedis()
	default:
		return errors.New("Invalid Database Selected. ")
	}
}

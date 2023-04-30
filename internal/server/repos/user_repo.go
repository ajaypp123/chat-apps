package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/kvstore"
	"github.com/ajaypp123/chat-apps/internal/server/models"
)

type UserRepoI interface {
	GetUser(ctx context.Context, username string) *models.Response
	PostUser(ctx context.Context, user *models.User) *models.Response
}

type UserRepo struct {
	// TODO: Create repository to store to database
	kv kvstore.KVStoreI
}

var userRepo *UserRepo = nil

func GetUserRepo() UserRepoI {
	if userRepo == nil {
		kvs, err := kvstore.NewRedisKVStore("localhost:6379", "", 0)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		userRepo = &UserRepo{
			kv: kvs,
		}
	}
	return userRepo
}

func (u *UserRepo) GetUser(ctx context.Context, username string) *models.Response {
	errRes := &models.Response{
		Code:   http.StatusNotFound,
		Status: models.Failed,
		Data:   "User not found",
	}
	uStr, err := u.kv.Get(username)
	if err != nil {
		fmt.Println(err)
		return errRes
	}
	fmt.Println(uStr)
	var usr models.User
	if err := json.Unmarshal([]byte(uStr), &usr); err != nil {
		fmt.Println(err)
		return errRes
	}
	return &models.Response{
		Code:   http.StatusOK,
		Status: models.Successed,
		Data:   &usr,
	}
}

func (u *UserRepo) PostUser(ctx context.Context, user *models.User) *models.Response {
	errRes := &models.Response{
		Code:   http.StatusNotFound,
		Status: models.Failed,
		Data:   "Failed",
	}

	user.Secret = common.GetUUID()
	uByte, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return errRes
	}

	if err := u.kv.Put(user.Username, string(uByte)); err != nil {
		fmt.Println(err)
		return errRes
	}
	return &models.Response{
		Code:   http.StatusOK,
		Status: models.Successed,
		Data:   user,
	}
}

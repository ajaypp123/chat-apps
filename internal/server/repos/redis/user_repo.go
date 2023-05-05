package redis

import (
	"encoding/json"
	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/common/logger"
	"github.com/ajaypp123/chat-apps/internal/server/repos"
	"net/http"

	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/kvstore"
	"github.com/ajaypp123/chat-apps/internal/server/models"
)

type UserRepo struct {
	// TODO: Add dao layer
	kv kvstore.KVStoreI
}

var userRepo *UserRepo = nil

func GetUserRepo(ctx *appcontext.AppContext) repos.UserRepoI {
	if userRepo == nil {
		kvs, err := kvstore.GetRedisKVStore()
		if err != nil {
			logger.Error(ctx, err)
			return nil
		}
		userRepo = &UserRepo{
			kv: kvs,
		}
	}
	return userRepo
}

func (u *UserRepo) GetUser(ctx *appcontext.AppContext, username string) *models.Response {
	logMsg := "UserRepo::GetUser username: " + username + " "
	logger.Info(ctx, "entry ", logMsg)
	errRes := &models.Response{
		Code:   http.StatusNotFound,
		Status: models.Failed,
		Data:   "User not found",
	}
	uStr, err := u.kv.Get(username)
	if err != nil {
		logger.Error(ctx, logMsg, err)
		return errRes
	}

	logger.Debug(ctx, logMsg, "User detail ", uStr)

	var usr models.User
	if err := json.Unmarshal([]byte(uStr), &usr); err != nil {
		logger.Error(ctx, logMsg, err)
		return errRes
	}
	logger.Info(ctx, "exit ", logMsg)
	return &models.Response{
		Code:   http.StatusOK,
		Status: models.Successed,
		Data:   &usr,
	}
}

func (u *UserRepo) PostUser(ctx *appcontext.AppContext, user *models.User) *models.Response {
	logMsg := "UserRepo::PostUser username: " + user.Username + " "
	logger.Info(ctx, "entry ", logMsg)
	errRes := &models.Response{
		Code:   http.StatusNotFound,
		Status: models.Failed,
		Data:   "Failed",
	}

	user.Secret = common.GetUUID()
	uByte, err := json.Marshal(user)
	if err != nil {
		logger.Error(ctx, logMsg, err)
		return errRes
	}

	if err := u.kv.Put(user.Username, string(uByte)); err != nil {
		logger.Error(ctx, logMsg, err)
		return errRes
	}
	logger.Info(ctx, "exit ", logMsg)
	return &models.Response{
		Code:   http.StatusOK,
		Status: models.Successed,
		Data:   user,
	}
}

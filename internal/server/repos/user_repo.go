package repos

import (
	"context"
	"net/http"

	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/internal/server/models"
)

type UserRepoI interface {
	GetUser(ctx context.Context, username string) *models.Response
	PostUser(ctx context.Context, username, name, phone string) *models.Response
}

type UserRepo struct {
	// TODO: Create repository to store to database
	usrMap map[string]*models.User
}

var userRepo *UserRepo = nil

func GetUserRepo() UserRepoI {
	if userRepo == nil {
		userRepo = &UserRepo{
			usrMap: make(map[string]*models.User),
		}
	}
	return userRepo
}

func (u *UserRepo) GetUser(ctx context.Context, username string) *models.Response {
	usr, ok := u.usrMap[username]
	if !ok {
		return &models.Response{
			Code:   http.StatusNotFound,
			Status: models.Failed,
			Data:   "User not found",
		}
	}
	return &models.Response{
		Code:   http.StatusOK,
		Status: models.Successed,
		Data:   usr,
	}
}

func (u *UserRepo) PostUser(ctx context.Context, username, name, phone string) *models.Response {
	usr := &models.User{
		Username: username,
		Secret:   common.GetUUID(),
		Name:     name,
		Phone:    phone,
	}

	if _, ok := u.usrMap[username]; ok {
		return &models.Response{
			Code:   http.StatusBadRequest,
			Status: models.Failed,
			Data:   "User Already exists",
		}
	}
	u.usrMap[usr.Username] = usr
	return &models.Response{
		Code:   http.StatusOK,
		Status: models.Successed,
		Data:   usr,
	}
}

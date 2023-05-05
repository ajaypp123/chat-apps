package services

import (
	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/internal/server/repos/redis"

	"github.com/ajaypp123/chat-apps/internal/server/models"
)

type UserServiceI interface {
	GetUser(ctx *appcontext.AppContext, username string) *models.Response
	PostUser(ctx *appcontext.AppContext, user *models.User) *models.Response
}

type UserService struct {
}

func NewUserService() UserServiceI {
	return &UserService{}
}

func (u *UserService) GetUser(ctx *appcontext.AppContext, username string) *models.Response {
	return redis.GetUserRepo(ctx).GetUser(ctx, username)
}

func (u *UserService) PostUser(ctx *appcontext.AppContext, user *models.User) *models.Response {
	return redis.GetUserRepo(ctx).PostUser(ctx, user)
}

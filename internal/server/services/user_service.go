package services

import (
	"context"

	"github.com/ajaypp123/chat-apps/internal/server/models"
	"github.com/ajaypp123/chat-apps/internal/server/repos"
)

type UserServiceI interface {
	GetUser(ctx context.Context, username string) *models.Response
	PostUser(ctx context.Context, user *models.User) *models.Response
}

type UserService struct {
}

func NewUserService() UserServiceI {
	return &UserService{}
}

func (u *UserService) GetUser(ctx context.Context, username string) *models.Response {
	return repos.GetUserRepo().GetUser(ctx, username)
}

func (u *UserService) PostUser(ctx context.Context, user *models.User) *models.Response {
	return repos.GetUserRepo().PostUser(ctx, user)
}

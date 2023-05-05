package repos

import (
	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/internal/server/models"
)

type UserRepoI interface {
	GetUser(ctx *appcontext.AppContext, username string) *models.Response
	PostUser(ctx *appcontext.AppContext, user *models.User) *models.Response
}

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ajaypp123/chat-apps/common/logger"
	"net/http"

	"github.com/ajaypp123/chat-apps/internal/server/models"
	"github.com/ajaypp123/chat-apps/internal/server/services"
	"github.com/gorilla/mux"
)

type UserController struct {
	userService services.UserServiceI
}

func NewUserController(userService services.UserServiceI) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) RegisterUserHandler(r *mux.Router) {
	// TODO: can we generate from json
	r.HandleFunc(getUsers, c.getUsers).Methods(http.MethodGet)
	r.HandleFunc(postUsers, c.createUser).Methods(http.MethodPost)
}

func (c *UserController) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := getCtx(r)
	username := r.URL.Query().Get("username")
	resp := c.userService.GetUser(ctx, username)

	encodeResponse(ctx, w, resp)
}

func (c *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := getCtx(r)
	var user *models.User
	if e := json.NewDecoder(r.Body).Decode(&user); e != nil {
		logger.Error(ctx, "Failed to parse request, Err: ", e)
		encodeError(ctx, w, invalidPayload, errors.New(fmt.Sprintf("Failed to parse request, Err: %v", e)))
		return
	}
	resp := c.userService.PostUser(ctx, user)

	encodeResponse(ctx, w, resp)
}

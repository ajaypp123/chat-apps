package controller

import (
	"context"
	"encoding/json"
	"fmt"
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

func (u *UserController) RegisterUserHandler(r *mux.Router) {
	// TODO: can we genarate from json
	r.HandleFunc("/v1/chatapp/users", u.getUsers).Methods("GET")
	r.HandleFunc("/v1/chatapp/users", u.createUser).Methods("POST")
}

func (c *UserController) getUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	resp := c.userService.GetUser(context.Background(), username)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.Code)
	json.NewEncoder(w).Encode(resp)
}

func (c *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if e := json.NewDecoder(r.Body).Decode(&user); e != nil {
		fmt.Println("failed")
		return
	}
	resp := c.userService.PostUser(context.Background(), user)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.Code)
	json.NewEncoder(w).Encode(resp)
}

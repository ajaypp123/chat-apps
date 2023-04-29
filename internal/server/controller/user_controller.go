package controller

import (
	"context"
	"encoding/json"
	"net/http"

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
	// TODO: take from body
	name := r.URL.Query().Get("name")
	phone := r.URL.Query().Get("phone")
	username := r.URL.Query().Get("username")
	resp := c.userService.PostUser(context.Background(), username, name, phone)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.Code)
	json.NewEncoder(w).Encode(resp)
}

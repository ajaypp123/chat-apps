package models

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Secret   string `json:"secret"`
}

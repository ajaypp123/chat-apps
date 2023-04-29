package models

type User struct {
	Username string `json:"Username"`
	Name     string `json:"Name"`
	Phone    string `json:"Phone"`
	Secret   string `json:"Secret"`
}

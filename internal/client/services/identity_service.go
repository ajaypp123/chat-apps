package services

type ClientData struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Secret   string `json:"secret"`
}

var Data ClientData

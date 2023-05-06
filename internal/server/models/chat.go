package models

type ConnectionDetail struct {
	Topic     string `json:"topic"`
	Partition string `json:"partition"`
	Username  string `json:"username"`
}

package models

const (
	Failed    = "failed"
	Successed = "success"
)

type Response struct {
	ReqId  string      `json:"req_id"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Code   int         `json:"code"`
}

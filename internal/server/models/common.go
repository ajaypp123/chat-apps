package models

const (
	Failed    = "failed"
	Successed = "success"
)

type Response struct {
	ReqId  string      `json:"req_id"`
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

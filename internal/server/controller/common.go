package controller

import (
	"encoding/json"
	"github.com/ajaypp123/chat-apps/common"
	"github.com/ajaypp123/chat-apps/common/appcontext"
	"github.com/ajaypp123/chat-apps/internal/server/models"
	"net/http"
)

var (
	getUsers   = "/v1/chat-apps/users"
	postUsers  = "/v1/chat-apps/users"
	getHealths = "/v1/chat-apps/health"
)

const (
	invalidPayload = iota
	badRouting
)

func getCtx(r *http.Request) *appcontext.AppContext {
	ctx := appcontext.DefaultContext.DeepCopy()
	ctx.AddValue("req-id", common.GetUUID())
	ctx.AddValue("url", r.URL)
	ctx.AddValue("method", r.Method)
	ctx.AddValue("addr", r.RemoteAddr)
	return ctx
}

func encodeResponse(_ *appcontext.AppContext, w http.ResponseWriter, resp *models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.Code)
	json.NewEncoder(w).Encode(resp)
}

func encodeError(ctx *appcontext.AppContext, w http.ResponseWriter, errCode int, err error) {
	resp := &models.Response{
		ReqId:  ctx.GetValue("req-id").(string),
		Status: models.Failed,
		Data:   err,
		Code:   0,
	}
	switch errCode {
	case invalidPayload:
		resp.Code = 400
	case badRouting:
		resp.Code = 404
	default:
		resp.Code = 400
	}
	encodeResponse(ctx, w, resp)
}

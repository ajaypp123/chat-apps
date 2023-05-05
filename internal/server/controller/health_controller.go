package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) RegisterHealthHandler(r *mux.Router) {
	r.HandleFunc(getHealths, h.getHealths).Methods(http.MethodGet)
}

func (h *HealthController) getHealths(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode([]byte(`OK`))
}

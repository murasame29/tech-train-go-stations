package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	respBody := &model.HealthzResponse{Message: fmt.Sprintf("OK OSname : %s", r.Context().Value(middleware.OsName{}).(string))}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		log.Println("json marshal error :", err)
		return
	}
}

package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	// healthzに関するルータを定義
	healthRouter(mux)

	return mux
}

// healthzに関するルータを定義
func healthRouter(mux *http.ServeMux) {
	healthz := handler.NewHealthzHandler()
	mux.HandleFunc("/healthz", healthz.ServeHTTP)

}

package router

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	// healthzに関するルータを定義
	healthRouter(mux)
	// todoに関するルータを定義
	todoRouter(mux, todoDB)
	// panicに関するルータを定義
	panicRouter(mux)

	return mux
}

// healthzに関するルータを定義
func healthRouter(mux *http.ServeMux) {
	healthz := handler.NewHealthzHandler()
	mux.Handle("/healthz", middleware.GetUserAgent(http.HandlerFunc(healthz.ServeHTTP)))

}

// todoに関するルータを定義
func todoRouter(mux *http.ServeMux, db *sql.DB) {
	todo := handler.NewTODOHandler(service.NewTODOService(db))

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.Method {
		case http.MethodGet:
			err = responseJson(todo.ReadTodo(w, r))
		case http.MethodPost:
			err = responseJson(todo.CreateTodo(w, r))
		case http.MethodPut:
			err = responseJson(todo.UpdateTodo(w, r))
		case http.MethodDelete:
			err = responseJson(todo.DeleteTodo(w, r))
		default:
			// TODO:エラーハンドリングする
		}
		//エラーがあった場合ログ出力して500を返す
		if err != nil {
			log.Println(err)
			responseJson(w, http.StatusInternalServerError, err)
		}
	})
}

func panicRouter(mux *http.ServeMux) {
	ph := handler.NewPanichandler()
	mux.Handle("/do-panic", middleware.Recovery(http.HandlerFunc(ph.ServeHTTP)))
}

// 任意のstatusをヘッドに入れたレスポンスを返す
func responseJson(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

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

type Router struct {
	TodoDB *sql.DB
	Mux    *http.ServeMux
}

func NewRouter(todoDB *sql.DB) *Router {
	// register routes
	mux := http.NewServeMux()

	router := &Router{
		Mux:    mux,
		TodoDB: todoDB,
	}
	router.healthRouter()
	router.panicRouter()
	router.todoRouter()

	return router
}

func (r *Router) healthRouter() {
	healthz := handler.NewHealthzHandler()
	r.Mux.Handle("/healthz", buildChain(
		http.HandlerFunc(healthz.ServeHTTP),
		middleware.Recovery,
		middleware.GetUserAgent,
		middleware.AccessLog,
	))
}

// todoに関するルータを定義
func (r *Router) todoRouter() {
	todo := handler.NewTODOHandler(service.NewTODOService(r.TodoDB))

	r.Mux.Handle("/todos", buildChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}),
		middleware.Recovery,
		middleware.GetUserAgent,
		middleware.AccessLog,
	))
}

func (r *Router) panicRouter() {
	ph := handler.NewPanichandler()
	r.Mux.Handle("/do-panic",
		buildChain(http.HandlerFunc(ph.ServeHTTP),
			middleware.Recovery,
			middleware.GetUserAgent,
			middleware.AccessLog,
		))
}

// 任意のstatusをヘッドに入れたレスポンスを返す
func responseJson(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}

func buildChain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](buildChain(h, m[1:cap(m)]...))
}

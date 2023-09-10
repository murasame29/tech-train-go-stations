package router

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/env"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/handler/response"
	"github.com/TechBowl-japan/go-stations/service"
)

type Router struct {
	TodoDB     *sql.DB
	middleware *middleware.Middleware
	Mux        *http.ServeMux
}

func NewRouter(todoDB *sql.DB, env *env.Env) *Router {
	// register routes
	mux := http.NewServeMux()

	router := &Router{
		Mux:        mux,
		middleware: middleware.NewMiddleware(env),
		TodoDB:     todoDB,
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
		r.middleware.Recovery,
		r.middleware.GetUserAgent,
		r.middleware.AccessLog,
	))
}

// todoに関するルータを定義
func (r *Router) todoRouter() {
	todo := handler.NewTODOHandler(service.NewTODOService(r.TodoDB))

	r.Mux.Handle("/todos", buildChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.Method {
		case http.MethodGet:
			err = response.ResponseJson(todo.ReadTodo(w, r))
		case http.MethodPost:
			err = response.ResponseJson(todo.CreateTodo(w, r))
		case http.MethodPut:
			err = response.ResponseJson(todo.UpdateTodo(w, r))
		case http.MethodDelete:
			err = response.ResponseJson(todo.DeleteTodo(w, r))
		default:
			// TODO:エラーハンドリングする
		}
		//エラーがあった場合ログ出力して500を返す
		if err != nil {
			log.Println(err)
			response.ResponseJson(w, http.StatusInternalServerError, err)
		}
	}),
		r.middleware.Recovery,
		r.middleware.GetUserAgent,
		r.middleware.AccessLog,
	))
}

func (r *Router) panicRouter() {
	ph := handler.NewPanichandler()
	r.Mux.Handle("/do-panic",
		buildChain(http.HandlerFunc(ph.ServeHTTP),
			r.middleware.Recovery,
			r.middleware.GetUserAgent,
			r.middleware.AccessLog,
		))
}

func buildChain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](buildChain(h, m[1:cap(m)]...))
}

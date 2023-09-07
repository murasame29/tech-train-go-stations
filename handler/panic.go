package handler

import (
	"net/http"
)

type PanicHandler struct{}

func NewPanichandler() http.Handler {
	return &PanicHandler{}
}

func (*PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic!!!")
}

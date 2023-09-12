package middleware

import (
	"net/http"

	"github.com/TechBowl-japan/go-stations/env"
	"github.com/TechBowl-japan/go-stations/handler/response"
	"github.com/TechBowl-japan/go-stations/model"
)

// https://pkg.go.dev/encoding/base64

func (m *Middleware) BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, password, ok := r.BasicAuth()
		if !ok {
			response.Unauthorized(w, model.ErrorResponse{Error: "UnAuthorized"})
			return
		}

		if !certification(userID, password, m.env) {
			response.Unauthorized(w, model.ErrorResponse{Error: "Wrong userId or password"})
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func certification(userID, password string, env *env.Env) bool {
	return env.UserID == userID && env.Password == password
}

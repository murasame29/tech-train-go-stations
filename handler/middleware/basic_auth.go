package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/TechBowl-japan/go-stations/env"
	"github.com/TechBowl-japan/go-stations/handler/response"
	"github.com/TechBowl-japan/go-stations/model"
)

// https://pkg.go.dev/encoding/base64

const (
	AuthenticateHeaderKey = "Authorization"
)

const (
	USER_ID = iota
	PASSWORD
)

func (m *Middleware) BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authToken := getHeader(r, AuthenticateHeaderKey)
		// authTokenが空の場合
		if len(authToken) == 0 { //Equal
			response.Unauthorized(w, model.ErrorResponse{Error: "Token cannot be empty"})
			return
		}

		rawToken, err := base64Decode(authToken)
		if err != nil {
			response.Unauthorized(w, model.ErrorResponse{Error: fmt.Sprintf("base64Decode Error : %s", err)})
			return
		}

		token := strings.Split(string(rawToken), ":")

		if !certification(token, m.env) {
			response.Unauthorized(w, model.ErrorResponse{Error: "Wrong userId or password"})
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func base64Decode(token string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(token)
}

func certification(token []string, env *env.Env) bool {
	return env.UserID == token[USER_ID] && env.Password == token[PASSWORD]
}

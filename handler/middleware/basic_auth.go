package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/TechBowl-japan/go-stations/env"
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
			// 401を返す
		}

		rawToken, err := base64Decode(authToken)
		if err != nil {
			// 返るErrorがCorruptInputErrorだから401で返す(base64のソース参照)
		}

		token := strings.Split(string(rawToken), ":")

		if !certification(token, m.env) {
			// UserID or Passwordの不一致 401
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

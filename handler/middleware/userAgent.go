package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

const (
	OsName = "os_name"
)

func GetUserAgent(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ユーザエージェントの取得
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), OsName, ua.OS)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

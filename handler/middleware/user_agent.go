package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type OsName struct{}

func getUserAgent(userAgent string) useragent.UserAgent {
	return useragent.Parse(userAgent)
}

func (m *Middleware) GetUserAgent(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ユーザエージェントの取得
		ua := getUserAgent(r.UserAgent())
		ctx := context.WithValue(r.Context(), OsName{}, ua.OS)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

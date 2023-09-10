package middleware

import (
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler/response"
)

// tips : 多重入れ子構造は可読性を損なう可能性あり
func (m *Middleware) Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				response.InternalServerError(w, nil)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

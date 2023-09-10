package middleware

import (
	"encoding/json"
	"net/http"
)

// tips : 多重入れ子構造は可読性を損なう可能性あり
func (m *Middleware) Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(nil)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

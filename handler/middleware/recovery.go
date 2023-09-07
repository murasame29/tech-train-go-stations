package middleware

import (
	"encoding/json"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(nil)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

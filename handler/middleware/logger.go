package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/handler/response"
)

type AccessLogger struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency(ms)"`
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

// 開始時間から何秒かかったか (ms)で出力
func getLatency(start time.Time) int64 {
	return int64(time.Since(start) / time.Microsecond)
}

func (m *Middleware) AccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var logger AccessLogger

		logger.Path = r.URL.Path
		logger.Timestamp = time.Now()
		logger.OS = (r.Context()).Value(OsName{}).(string)

		defer func() {
			logger.Latency = getLatency(logger.Timestamp)
			body, err := json.Marshal(logger)
			if err != nil {
				response.InternalServerError(w, err)
			}
			// 標準出力
			fmt.Println(string(body))
		}()
		h.ServeHTTP(w, r)
	})
}

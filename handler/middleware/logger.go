package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"`
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

// 開始時間から何秒かかったか (ms)で出力
func getLatency(start time.Time) int64 {
	return int64(time.Since(start) / time.Microsecond)
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var log Log

		log.Path = r.URL.Path
		log.Timestamp = time.Now()
		// 最後に行われる処理
		defer func() {
			log.OS = r.Context().Value(OsName{}).(string)
			log.Latency = getLatency(log.Timestamp)

			body, err := json.Marshal(log)
			if err != nil {
				if rec := recover(); rec != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(nil)
				}
			}
			// 標準出力
			fmt.Println(string(body))
		}()
	})
}

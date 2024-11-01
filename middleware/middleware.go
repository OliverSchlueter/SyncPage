package middleware

import (
	"net/http"
	"syncpage/logger"
	"time"
)

func Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		sr := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(sr, r)

		elapsedTime := time.Since(startTime)

		logger.InfoProps("Request received", map[string]interface{}{
			"method":       r.Method,
			"url":          r.URL.String(),
			"status":       sr.Status,
			"elapsed_time": elapsedTime.Milliseconds(),
			"user_agent":   r.UserAgent(),
			"referer":      r.Referer(),
		})
	}
}

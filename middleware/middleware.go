package middleware

import (
	"fmt"
	"net/http"
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
		fmt.Printf("[Req] %d %s: %s %s in %s\n", sr.Status, http.StatusText(sr.Status), r.Method, r.URL, elapsedTime)
	}
}

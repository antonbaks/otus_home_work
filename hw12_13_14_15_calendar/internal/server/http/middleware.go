package internalhttp

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)
		s := []string{
			r.RemoteAddr,
			"[" + time.Now().Format(time.RFC3339) + "]",
			r.Method,
			r.RequestURI,
			r.Proto,
			strconv.Itoa(recorder.Status),
			duration.String(),
			r.UserAgent(),
			"\n",
		}

		os.Stdout.WriteString(strings.Join(s, " "))
	})
}

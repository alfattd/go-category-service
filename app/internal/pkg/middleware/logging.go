package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/alfattd/category-service/internal/pkg/requestid"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			log.Info("request completed",
				"request_id", requestid.FromContext(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.status,
				"duration", time.Since(start).String(),
			)
		})
	}
}

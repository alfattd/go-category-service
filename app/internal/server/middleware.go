package server

import (
	"log/slog"
	"net/http"

	"github.com/alfattd/category-service/internal/pkg/monitor"
	"github.com/alfattd/category-service/internal/pkg/requestid"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			monitor.HttpRequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()
		}
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("incoming request",
			"request_id", requestid.FromContext(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
		)
		next.ServeHTTP(w, r)
	})
}

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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := requestid.FromContext(r.Context())
		slog.Info("incoming request",
			"request_id", rid,
			"method", r.Method,
			"path", r.URL.Path,
		)
		next.ServeHTTP(w, r)
	})
}

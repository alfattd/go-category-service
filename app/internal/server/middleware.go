package server

import (
	"net/http"

	"github.com/alfattd/category-service/internal/platform/monitor"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			monitor.HttpRequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()
		}
		next.ServeHTTP(w, r)
	})
}

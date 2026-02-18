package middleware

import (
	"net/http"

	"github.com/alfattd/category-service/internal/pkg/requestid"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}

		ctx := requestid.WithContext(r.Context(), id)
		w.Header().Set(RequestIDHeader, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package requestid

import "context"

type contextKey string

const RequestIDKey contextKey = "request_id"

func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

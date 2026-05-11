package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// UserIDFromContext extracts the authenticated user's UUID from the request
// context. Panics if called outside the Authenticate middleware — that is an
// invariant violation, not a runtime error to recover from.
func UserIDFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(ctxKeyUserID).(uuid.UUID)
	if !ok {
		panic("middleware.UserIDFromContext: user ID not found in context — route not protected by Authenticate middleware")
	}
	return id
}

// RequestID returns a per-request ID injected by the RequestID middleware.
// Returns empty string when middleware is not applied.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKeyRequestID).(string)
	return id
}

// SetRequestID is an http.Handler adapter that generates a unique request ID.
// The ID is added to both the context and the X-Request-ID response header.
func SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

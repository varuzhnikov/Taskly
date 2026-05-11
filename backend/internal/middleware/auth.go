package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/todoist/backend/internal/service"
)

type ctxKey string

const ctxKeyUserID ctxKey = "user_id"

// Authenticate validates the Bearer JWT and injects the parsed user ID into
// the request context. Requests without a valid token receive 401.
func Authenticate(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing or malformed authorization header"}`, http.StatusUnauthorized)
				return
			}

			userID, err := auth.ValidateAccessToken(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUserID, userID)))
		})
	}
}

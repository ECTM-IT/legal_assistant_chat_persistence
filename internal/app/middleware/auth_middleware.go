package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"go.uber.org/zap"
)

// contextKey is a type used for context keys to avoid key collisions.
type contextKey string

const uidContextKey contextKey = "uid"

// FirebaseAuthMiddleware is a middleware for authentication using Firebase.
func FirebaseAuthMiddleware(firebaseApp *firebase.App, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				logger.Warn("Authorization header is missing")
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}

			idToken := strings.TrimPrefix(authHeader, "Bearer ")
			client, err := firebaseApp.Auth(r.Context())
			if err != nil {
				logger.Error("Error getting Auth client", zap.Error(err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			token, err := client.VerifyIDToken(r.Context(), idToken)
			if err != nil {
				logger.Warn("Error verifying ID token", zap.Error(err))
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add the UID to the request context using the custom key type
			ctx := context.WithValue(r.Context(), uidContextKey, token.UID)
			logger.Info("Authenticated request", zap.String("uid", token.UID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

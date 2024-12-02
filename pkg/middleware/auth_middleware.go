package middleware

import (
	"cloudservices/pkg/auth"
	"context"
	"errors"
	"net/http"
	"strings"
)

// TokenAuthMiddleware enforces JWT-based authentication
func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Bearer token parsing
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate the token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Store claims in the request context
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireUsernameContext retrieves the username from the context
func RequireUsernameContext(r *http.Request) (string, error) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		return "", errors.New("username not found in context")
	}
	return username, nil
}

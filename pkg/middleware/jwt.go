package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"
)

// JWT is a middleware that authenticates the incoming request.
// This middleware parses the claims form the incoming JWT in the `Authorization` header and writes them to the request context.
// If the request is not authenticated, it returns a 401 Unauthorized response.
// If the request is authenticated, it calls the next handler in the chain.
func JWT(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Register exceptional routes.
			if r.URL.Path == "/health" ||
				r.URL.Path == "/metrics" ||
				r.URL.Path == "/signin" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract the JWT from the `Authorization` header.
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Remove the `Bearer ` prefix from the JWT.
			header = header[7:]

			// Parse the JWT and extract the claims.
			token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			})

			if err != nil {
				http.Error(w, fmt.Sprintf("failed to parse the JWT: %s", err), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "supplied JWT is invalid", http.StatusUnauthorized)
				return
			}

			// Parse the claims from the token.
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "failed to parse the claims", http.StatusUnauthorized)
				return
			}

			// Extract the user ID from the claims.
			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "failed to parse the user ID", http.StatusUnauthorized)
				return
			}

			// Write the user_id to the request context.
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserID, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// XJWTClaims is the key used to store the claims of the JWT in the context.
//
// The claims are used to store the information about the authenticated user.
const XJWTClaims Key = "x-jwt-claims"

type JWTClaims struct {
	jwt.StandardClaims
	XUserID uuid.UUID `json:"x-user-id"`
}

func (c JWTClaims) Valid() error {
	if c.XUserID == uuid.Nil {
		return fmt.Errorf("invalid user id")
	}
	return nil
}

//	JWT is a middleware that can be used to validate the JWTs.
//
// Generate temporary JWTs for testing from here: https://oauth.tools/collection/1712706959493-UZt
type JWTConfig struct {

	// Prefix is the type of the JWT.
	// Default: `Bearer`
	//
	// This field is optional.
	Prefix string

	// Algorithm is the algorithm of the key that will be used to validate the JWT.
	// Default: `HS256`
	//
	// This field is optional.
	Algorithm string

	// Issuer is the issuer of the JWT.
	// Default: ``
	//
	// This field is optional.
	Issuer string

	// Audience is the audience of the JWT.
	// Default: ``
	//
	// This field is optional.
	Audience string

	// Key is the secret key that will be used to validate the JWT.
	//
	// This field is mandatory.
	Key string

	// ExceptionalRoutes is the list of routes that will be excluded from the JWT validation.
	// For example, you can exclude the login route from the JWT validation.
	//
	// Example: []string{
	// 		"/login"
	// 		"/healthz"
	//	}
	//
	// This field is optional.
	ExceptionalRoutes []string

	// Header is the request header that will be used to extract the JWT from.
	// Default: `Authorization`
	//
	// This field is optional.
	Header string
}

func JWT(config *JWTConfig) Middleware {

	// Validate the configuration.
	if config == nil {
		panic("failed to initialize the JWT middleware: missing configuration")
	}

	if config.Key == "" {
		panic("failed to initialize the JWT middleware: missing key")
	}

	//
	// Set default values.
	//

	if config.Prefix == "" {
		config.Prefix = "Bearer"
	}

	if config.Algorithm == "" {
		config.Algorithm = "HS256"
	}

	if config.Header == "" {
		config.Header = "Authorization"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Avoid the JWT validation for the exceptional routes.
			for _, item := range config.ExceptionalRoutes {
				if r.URL.Path == item {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Extract the JWT from the appropriate header.
			header := r.Header.Get(config.Header)
			if header == "" {
				http.Error(w, "failed to extract the JWT from appropriate header", http.StatusUnauthorized)
				return
			}

			// Remove the prefix from the JWT.
			if len(header) > len(config.Prefix) && header[:len(config.Prefix)] == config.Prefix {
				header = header[len(config.Prefix)+1:]
			}

			// Parse the JWT and extract the claims.
			var claims JWTClaims
			token, err := jwt.ParseWithClaims(header, &claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Key), nil
			})

			if err != nil {
				http.Error(w, fmt.Sprintf("failed to parse the JWT: %s", err), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "supplied JWT is invalid", http.StatusUnauthorized)
				return
			}

			// Write the claims to the request context.
			r = r.WithContext(context.WithValue(r.Context(), XJWTClaims, claims))

			next.ServeHTTP(w, r)
		})
	}
}

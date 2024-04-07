package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type Key string

// X-Request-ID is the key used to store the request ID in the context and the response header.
const XRequestID Key = "X-Request-ID"

// RequestID middleware adds a unique UUID to the request context and response headers.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New().String()

		// Add the request ID to the request context.
		ctx = context.WithValue(ctx, XRequestID, id)

		// Update the request with the new context.
		r = r.WithContext(ctx)

		// Add the request ID to the response headers.
		w.Header().Set(string(XRequestID), id)
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

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

// TraceID middleware adds a unique UUID to the request context and response headers.
func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New().String()

		// Add the trace ID to the request context.
		ctx = context.WithValue(ctx, XTraceID, id)

		// Update the request with the new context.
		r = r.WithContext(ctx)

		// Add the trace ID to the response headers.
		w.Header().Set(string(XTraceID), id)
		next.ServeHTTP(w, r)
	})
}

// CorrelationID middleware adds a unique UUID to the request context and response headers.
func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New().String()

		// Add the correlation ID to the request context.
		ctx = context.WithValue(ctx, XCorrelationID, id)

		// Update the request with the new context.
		r = r.WithContext(ctx)

		// Add the correlation ID to the response headers.
		w.Header().Set(string(XCorrelationID), id)
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type Key string

// X-Request-ID is the key used to store the request ID in the context and the response header.
//
// The request ID is used to uniquely identify the request.
const XRequestID Key = "X-Request-ID"

// X-Trace-ID is the key used to store the trace ID in the context and the response header.
//
// The trace ID is used to trace the request through multiple services.
const XTraceID Key = "X-Trace-ID"

// X-Correlation-ID is the key used to store the correlation ID in the context and the response header.
//
// The correlation ID is used to correlate the request with other requests.
const XCorrelationID Key = "X-Correlation-ID"

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

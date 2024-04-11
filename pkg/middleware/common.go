package middleware

import (
	"net/http"
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

// X-JWT-Claims is the key used to store the claims of the JWT in the context.
//
// The claims are used to store the information about the authenticated user.
const XJWTClaims Key = "X-JWT-Claims"

type Middleware func(http.Handler) http.Handler

// Chain is a variadic function that executes multiple middlewares in sequential order.
func Chain(middlewares ...Middleware) Middleware {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

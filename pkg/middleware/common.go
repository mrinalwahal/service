package middleware

import (
	"net/http"
)

type Key string

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

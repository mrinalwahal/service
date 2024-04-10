package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Chain is a variadic function that executes multiple middlewares in sequential order.
func Chain[T func(http.Handler) http.Handler](middlewares ...T) T {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

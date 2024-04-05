package middleware

import (
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

// Logging middleware logs the incoming HTTP request.
// It logs the request method, request URL, and the time it took to process the request.
// It also logs the response status and response time.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Measure how long it takes to process the request
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info("%s %s", r.Method, r.URL.Path)
		slog.Info("request processed in %s", time.Since(start))
	})
}

package middleware

import (
	"log/slog"
	"net/http"
)

type RecoverConfig struct {

	// Logger is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	Logger *slog.Logger
}

// Recover is a middleware that recovers from the panics.
func Recover(config *RecoverConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if err == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(err)
					}

					config.Logger.LogAttrs(r.Context(), slog.LevelError, "panic recovered", slog.Attr{
						Key:   "error",
						Value: slog.AnyValue(err),
					})

					if r.Header.Get("Connection") != "Upgrade" {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

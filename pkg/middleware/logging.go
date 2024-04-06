package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			//writer := writer.Writer{ResponseWriter: w}
			next.ServeHTTP(w, r)

			attributes := []slog.Attr{
				{Key: "request_id", Value: slog.StringValue(r.Context().Value(XRequestID).(string))},
				// {Key: "status", Value: slog.IntValue(r.Status)},
				{Key: "duration", Value: slog.DurationValue(time.Since(start))},
				{Key: "hostname", Value: slog.StringValue(r.Host)},
				{Key: "method", Value: slog.StringValue(r.Method)},
				{Key: "path", Value: slog.StringValue(r.URL.Path)},
			}

			log.LogAttrs(r.Context(), slog.LevelInfo, "http request", attributes...)
		})
	}
}

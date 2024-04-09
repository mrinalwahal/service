package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrinalwahal/service/pkg/writer"
)

func Logging(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			//
			// If you want to run some code before the request is handled, you can do it here.
			// For example, you can  modify the request object before passing it to the next handler.
			// Like we do it in the `RequestID` middleware.
			//

			writer := writer.NewWriter(w)
			next.ServeHTTP(writer, r)

			//
			// If you want to run some code after the request is handled, you can do it here.
			// For our use case, we are going to log the request.
			//

			attributes := []slog.Attr{
				{Key: "request_id", Value: slog.StringValue(r.Context().Value(XRequestID).(string))},
				{Key: "status", Value: slog.IntValue(writer.Status())},
				{Key: "duration", Value: slog.DurationValue(time.Since(start))},
				{Key: "hostname", Value: slog.StringValue(r.Host)},
				{Key: "method", Value: slog.StringValue(r.Method)},
				{Key: "path", Value: slog.StringValue(r.URL.Path)},
			}

			log.LogAttrs(r.Context(), slog.LevelInfo, fmt.Sprintf("incoming %s request to %s", r.Method, r.URL.Path), attributes...)
		})
	}
}

package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrinalwahal/service/pkg/writer"
)

type LoggingConfig struct {

	// Logger is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	Logger *slog.Logger

	// LogLatency is the flag that determines if the latency of the request should be logged.
	// Latency is calculated as the difference between the time the request is received and the time the response is sent.
	// Default: `false`
	//
	// This field is optional.
	LogLatency bool

	// LogError is the flag that determines if the response status is 5xx then the error message should be logged.
	// Default: `false`
	//
	// This field is optional.
	LogError bool
}

func Logging(config *LoggingConfig) Middleware {

	// Set the default configuration.
	if config == nil {
		config = &LoggingConfig{}
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

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
				{Key: "timestamp", Value: slog.StringValue(start.String())},
				{Key: "request_id", Value: slog.StringValue(r.Context().Value(XRequestID).(string))},
				{Key: "status", Value: slog.IntValue(writer.Status())},
				{Key: "hostname", Value: slog.StringValue(r.Host)},
				{Key: "method", Value: slog.StringValue(r.Method)},
				{Key: "path", Value: slog.StringValue(r.URL.Path)},
			}

			if config.LogLatency {
				attributes = append(attributes, slog.Attr{Key: "latency", Value: slog.DurationValue(time.Since(start))})
			}

			// If the response status code is 5xx, log the error message.
			if writer.Status() >= 500 && config.LogError {

				// Parse the response data.
				// attributes = append(attributes, slog.Attr{Key: "error", Value: slog.StringValue(writer.Error())})

			} else {
				config.Logger.LogAttrs(r.Context(), slog.LevelInfo, fmt.Sprintf("incoming %s request to %s", r.Method, r.URL.Path), attributes...)
			}
		})
	}
}

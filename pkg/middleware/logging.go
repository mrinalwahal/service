package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Record is the primary structure in which a log record is maintained.
type Record struct {
	RequestID string        `json:"request_id,omitempty"`
	StartTime string        `json:"start_time,omitempty"`
	Status    int           `json:"status,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
	Hostname  string        `json:"hostname,omitempty"`
	Method    string        `json:"method,omitempty"`
	Path      string        `json:"path,omitempty"`
}

func (r *Record) attributes() []slog.Attr {

	// If the status is 0, then set it to 200.
	if r.Status == 0 {
		r.Status = http.StatusOK
	}

	return []slog.Attr{
		{Key: "request_id", Value: slog.StringValue(r.RequestID)},
		{Key: "start_time", Value: slog.StringValue(r.StartTime)},
		{Key: "status", Value: slog.IntValue(r.Status)},
		{Key: "duration", Value: slog.DurationValue(r.Duration)},
		{Key: "hostname", Value: slog.StringValue(r.Hostname)},
		{Key: "method", Value: slog.StringValue(r.Method)},
		{Key: "path", Value: slog.StringValue(r.Path)},
	}
}

func Logging(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			//writer := writer.Writer{ResponseWriter: w}
			next.ServeHTTP(w, r)

			record := &Record{
				RequestID: r.Context().Value(XRequestID).(string),
				StartTime: start.Format(time.RFC3339),
				//Status:    writer.Status(),
				Duration: time.Since(start),
				Hostname: r.Host,
				Method:   r.Method,
				Path:     r.URL.Path,
			}

			log.LogAttrs(r.Context(), slog.LevelInfo, "request processed", record.attributes()...)
		})
	}
}

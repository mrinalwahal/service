package logger

import (
	"log/slog"
	"os"
	"time"
)

// Log is the primary structure in which a log record is maintained.
type Log struct {
	Method       string
	Path         string
	Status       int
	ResponseTime time.Duration
}

func New() *Log {
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
	l := slog.New(h)
	l.WithGroup("service")

	return &Log{}
}

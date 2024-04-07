package handler

import (
	"context"
	"net/http"
)

// Handler is the interface that wraps the basic methods of a handler.
type Handler interface {
	Validate(context.Context) error
	Process(context.Context) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type HTTPHandlerConfig struct {

	//	Logger is the `log/slog` instance that will be used to log messages.
	//	Default: `slog.DefaultLogger`
	//
	//	This field is optional.
	Logger *slog.Logger
}

type HTTPHandler struct {
	log *slog.Logger
}

// NewHTTPHandler creates a new instance of `HTTPHandler`.
func NewHTTPHandler(config *HTTPHandlerConfig) *HTTPHandler {
	if config.Logger == nil {
		config.Logger = slog.Default()
	}
	handler := HTTPHandler{
		log: config.Logger,
	}

	return &handler
}

// Handler is a function that handles incoming requests.
// type Handler func(*http.Request) error

//
// Functions which will handle incoming requests.
//

// Create handler create a new record.
func (h *HTTPHandler) Create(req *http.Request) error {

	// Prepare the context from the request context.
	return &Response{
		Status: http.StatusNotImplemented,
	}
}

// Get handler retrieves a specific record by it's UUID.
func (h *HTTPHandler) Get(req *http.Request) error {

	// Get the record's UUID from the request path.
	_, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid record ID.",
			Err:     err,
		}
	}

	return &Response{
		Status: http.StatusNotImplemented,
	}
}

// // List handler retrieves all records.
// func (h *HTTPHandler) List(req *http.Request) error {
// 	return WriteString(w, http.StatusNotImplemented, "Not implemented")
// }

// // Update handler updates a specific record by it's UUID.
// func (h *HTTPHandler) Update(req *http.Request) error {
// 	return WriteString(w, http.StatusNotImplemented, "Not implemented")
// }

// // Delete handler deletes a specific record by it's UUID.
// func (h *HTTPHandler) Delete(req *http.Request) error {
// 	return WriteString(w, http.StatusNotImplemented, "Not implemented")
// }

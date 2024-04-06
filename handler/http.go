package handler

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

type HTTPHandlerConfig struct {

	// Database dialector to use.
	// Default: `gorm.PostgresDialect`
	//
	// This field is mandatory.
	Dialector *gorm.Dialector

	// Prefix is the prefix for all the routes.
	// Example: `/v1`
	// Default: `""`
	//
	// This field is optional.
	Prefix string

	// Router is a `http.ServeMux` instance that will be used to register the routes.
	// If this field is provided, all the handlers with default route patterns will be automatically registered on this router.
	// Avoid using this field if you want to register the routes manually.
	// Default: `nil`
	//
	// This field is optional.
	// Router *http.ServeMux

	//	Logger is the `log/slog` instance that will be used to log messages.
	//	Default: `slog.DefaultLogger`
	//
	//	This field is optional.
	Logger *slog.Logger
}

type HTTPHandler struct {
	dialector *gorm.Dialector
	log       *slog.Logger
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the handler function for the request.
	// handlerFunc, ok := h.routes[r.Method+" "+r.URL.Path]
}

// NewHTTPHandler creates a new instance of `HTTPHandler`.
func NewHTTPHandler(config *HTTPHandlerConfig) *HTTPHandler {

	handler := HTTPHandler{
		log:       config.Logger,
		dialector: config.Dialector,
	}

	// Set the default logger if not provided.
	if handler.log == nil {
		handler.log = slog.Default()
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

// // Get handler retrieves a specific record by it's UUID.
// func (h *HTTPHandler) Get(req *http.Request) error {

// 	// Get the record's UUID from the request path.
// 	_, err := uuid.Parse(req.PathValue("id"))
// 	if err != nil {
// 		return &Response{
// 			Status:  http.StatusBadRequest,
// 			Message: "Invalid record ID.",
// 			Err:     err,
// 		}
// 	}

// 	return &Response{
// 		Status: http.StatusNotImplemented,
// 	}
// }

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

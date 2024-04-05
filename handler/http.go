package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type HTTPHandlerConfig struct {

	//	Prefix is the prefix that will be added to all the routes.
	//	Example: "/api" or "/v1" or "/api/v1".
	//
	//	This field is optional.
	Prefix string

	//	Logger is the `log/slog` instance that will be used to log messages.
	//	Default: `slog.DefaultLogger`
	//
	//	This field is optional.
	Logger *slog.Logger
}

type HTTPHandler struct {
	*http.ServeMux
	prefix string
	log    *slog.Logger
}

// NewHTTPHandler creates a new instance of `HTTPHandler`.
func NewHTTPHandler(config *HTTPHandlerConfig) *HTTPHandler {
	if config.Logger == nil {
		config.Logger = slog.Default()
	}
	return &HTTPHandler{
		ServeMux: http.NewServeMux(),
		log:      config.Logger,
	}
}

// Register default routes.
func (r *HTTPHandler) Register() {

	// Health check endpoint.
	// Returns OK if the server is running.
	r.HandleFunc("GET /healthz", func(r *http.Request) error {
		return &Response{
			Status: http.StatusOK,
		}
	})

	// CRUD routes.
	r.HandleFunc("POST /", r.Create)
	r.HandleFunc("GET /{id}", r.Get)
	// r.HandleFunc("GET /", r.List)
	// r.HandleFunc("PUT /{id}", r.Update)
	// r.HandleFunc("DELETE /{id}", r.Delete)
}

// HandleFunc registers a new route with the router.
func (r *HTTPHandler) HandleFunc(pattern string, handler func(*http.Request) error) {
	r.ServeMux.HandleFunc(r.prefix+pattern, func(w http.ResponseWriter, req *http.Request) {
		r.log.Info("Request received", slog.String("method", req.Method), slog.String("path", req.URL.Path))
		if err := handler(req); err != nil {
			r.log.Error("Request failed", slog.String("method", req.Method), slog.String("path", req.URL.Path), err)

			// Run type assertion on the response to check if it is of type `Response`.
			// If it is, then write the response as JSON.
			// If it is not, then wrap the error in a new `Response` structure with defaults.
			if response, ok := err.(*Response); ok {
				WriteJSON(w, response.Status, response)
				return
			}
			WriteJSON(w, http.StatusInternalServerError, &Response{
				Message: "Your broke something on our server :(",
				Err:     err,
			})
			return
		}
	})
}

// Functions which will handle incoming requests.
//
// Create handler create a new record.
func (r *HTTPHandler) Create(req *http.Request) error {

	// Prepare the context from the request context.
	return &Response{
		Status: http.StatusNotImplemented,
	}
}

// Get handler retrieves a specific record by it's UUID.
func (r *HTTPHandler) Get(req *http.Request) error {

	// Get the record's UUID from the request path.
	_, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	return &Response{
		Status: http.StatusNotImplemented,
	}
}

// // List handler retrieves all records.
// func (r *HTTPHandler) List(req *http.Request) error {
// 	return WriteString(w, http.StatusNotImplemented, "Not implemented")
// }

// // Update handler updates a specific record by it's UUID.
// func (r *HTTPHandler) Update(req *http.Request) error {
// 	return WriteString(w, http.StatusNotImplemented, "Not implemented")
// }

// // Delete handler deletes a specific record by it's UUID.
// func (r *HTTPHandler) Delete(req *http.Request) error {
// 	return WriteString(w, http.StatusNotImplemented, "Not implemented")
// }

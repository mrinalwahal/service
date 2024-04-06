package router

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type HTTPRouter struct {
	*http.ServeMux

	// Logger is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	Logger *slog.Logger
}

type HTTPRouterConfig struct {

	// Logger is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	Logger *slog.Logger
}

// NewHTTPRouter creates a new instance of `HTTPRouter`.
func NewHTTPRouter(config *HTTPRouterConfig) *HTTPRouter {

	router := HTTPRouter{
		ServeMux: http.NewServeMux(),
		Logger:   config.Logger,
	}

	// Set the default logger if not provided.
	if router.Logger == nil {
		router.Logger = slog.Default()
	}

	// Register the default routes.
	router.registerDefaultRoutes()

	return &router
}

// registerDefaultRoutes registers the default routes.
func (r *HTTPRouter) registerDefaultRoutes() {
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	r.HandleFunc("POST /", handle(r.Create))
}

// HandleFunc registers the handler function for the given pattern.
// func (r *HTTPRouter) HandleFunc(pattern string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {

// }

// ServeHTTP handles the incoming HTTP request.
// func (r *HTTPRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {

// }

//
// Functions which will handle incoming requests.
//

// Create handler create a new record.
func (r *HTTPRouter) Create(req *http.Request) error {
	return fmt.Errorf("not implemented")

	// Decode the request body.
	body, err := decode[CreateOptions](req)
	if err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Err:     err,
		}
	}

	// Log the request body.
	r.Logger.LogAttrs(req.Context(), slog.LevelInfo, "create request", slog.String("title", body.Title))

	// Prepare the context from the request context.
	return &Response{
		Status:  http.StatusCreated,
		Message: "Record created successfully",
		Data:    body,
	}
}

//
// Utility functions.
//

func handle(handlerFunc func(*http.Request) error) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := handlerFunc(req); err != nil {

			// Run type assertion on the response to check if it is of type `Response`.
			// If it is, then write the response as JSON.
			// If it is not, then wrap the error in a new `Response` structure with defaults.
			if response, ok := err.(*Response); ok {
				write(w, response.Status, response)
				return
			}
			write(w, http.StatusInternalServerError, &Response{
				Message: "Your broke something on our server :(",
				Err:     err,
			})
			return
		}
	}
}

// write writes the data to the supplied http response writer.
func write(w http.ResponseWriter, status int, response any) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

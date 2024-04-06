package router

import (
	"log/slog"
	"net/http"
)

type HTTPRouter struct {

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
		Logger: config.Logger,
	}

	// Set the default logger if not provided.
	if router.Logger == nil {
		router.Logger = slog.Default()
	}

	// Register the default routes.
	router.HandleFunc("POST /", router.Create)

	return &router
}

// HandleFunc registers the handler function for the given pattern.
func (r *HTTPRouter) HandleFunc(pattern string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {

}

// ServeHTTP handles the incoming HTTP request.
func (r *HTTPRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

// ListenAndServe starts the HTTP server.
func (r *HTTPRouter) ListenAndServe(addr string) error {
	return nil
}

// ListenAndServeTLS starts the HTTPS server.
func (r *HTTPRouter) ListenAndServeTLS(addr, certFile, keyFile string) error {
	return nil
}

//
// Functions which will handle incoming requests.
//

// Create handler create a new record.
func (r *HTTPRouter) Create(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created"))
}

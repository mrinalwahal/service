package router

import (
	"log/slog"
	"net/http"

	"github.com/mrinalwahal/service/router/http/handlers"
	"gorm.io/gorm"
)

type HTTPRouter struct {
	*http.ServeMux

	// Database connection.
	// The connection should already be open.
	//
	// This field is mandatory.
	db *gorm.DB

	// log is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	log *slog.Logger
}

// HandleFunc registers the handler function for the given pattern.
// func (r *HTTPRouter) HandleFunc(pattern string, handlerFunc func(w http.ResponseWriter, req *http.Request)) {}

// ServeHTTP handles the incoming HTTP request.
// func (r *HTTPRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

type HTTPRouterConfig struct {

	// Database connection.
	// The connection should already be open.
	//
	// This field is mandatory.
	DB *gorm.DB

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
		db:       config.DB,
		log:      config.Logger,
	}

	// Set the default logger if not provided.
	if router.log == nil {
		router.log = slog.Default()
	}

	router.log = router.log.With("layer", "router")

	// Register the default routes.
	router.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Handle("POST /", handlers.NewCreateHandler(&handlers.CreateHandlerConfig{
		DB:     router.db,
		Logger: router.log,
	}))

	router.Handle("GET /", handlers.NewListHandler(&handlers.ListHandlerConfig{
		DB:     router.db,
		Logger: router.log,
	}))

	router.Handle("GET /{id}", handlers.NewGetHandler(&handlers.GetHandlerConfig{
		DB:     router.db,
		Logger: router.log,
	}))

	router.Handle("PATCH /{id}", handlers.NewUpdateHandler(&handlers.UpdateHandlerConfig{
		DB:     router.db,
		Logger: router.log,
	}))

	router.Handle("DELETE /{id}", handlers.NewDeleteHandler(&handlers.DeleteHandlerConfig{
		DB:     router.db,
		Logger: router.log,
	}))

	return &router
}

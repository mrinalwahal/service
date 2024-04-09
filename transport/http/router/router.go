package router

import (
	"log/slog"
	"net/http"

	"github.com/mrinalwahal/service/service"
	"github.com/mrinalwahal/service/transport/http/handler"
)

type HTTPRouter struct {
	*http.ServeMux

	// Service layer.
	//
	// This field is mandatory.
	service service.Service

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

	// Service layer.
	//
	// This field is mandatory.
	Service service.Service

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
		service:  config.Service,
		log:      config.Logger,
	}

	// Set the default logger if not provided.
	if router.log == nil {
		router.log = slog.Default()
	}

	// router.log = router.log.With("layer", "http")

	// Register the default routes.
	router.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Handle("POST /v1", handler.NewCreateHandler(&handler.CreateHandlerConfig{
		Service: config.Service,
		Logger:  router.log,
	}))

	router.Handle("GET /v1", handler.NewListHandler(&handler.ListHandlerConfig{
		Service: config.Service,
		Logger:  router.log,
	}))

	router.Handle("GET /v1/{id}", handler.NewGetHandler(&handler.GetHandlerConfig{
		Service: config.Service,
		Logger:  router.log,
	}))

	router.Handle("PATCH /v1/{id}", handler.NewUpdateHandler(&handler.UpdateHandlerConfig{
		Service: config.Service,
		Logger:  router.log,
	}))

	router.Handle("DELETE /v1/{id}", handler.NewDeleteHandler(&handler.DeleteHandlerConfig{
		Service: config.Service,
		Logger:  router.log,
	}))

	return &router
}

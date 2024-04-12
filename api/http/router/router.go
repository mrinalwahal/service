package router

import (
	"log/slog"
	"net/http"

	v1 "github.com/mrinalwahal/service/api/http/handlers/v1"
	"github.com/mrinalwahal/service/service"
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

	// Register the v1 routes.
	router.RegisterV1Routes()

	return &router
}

// RegisterV1Routes registers /v1 routes.
func (r *HTTPRouter) RegisterV1Routes() {

	r.Handle("POST /v1", v1.NewCreateHandler(&v1.CreateHandlerConfig{
		Service: r.service,
		Logger:  r.log,
	}))

	r.Handle("GET /v1", v1.NewListHandler(&v1.ListHandlerConfig{
		Service: r.service,
		Logger:  r.log,
	}))

	r.Handle("GET /v1/{id}", v1.NewGetHandler(&v1.GetHandlerConfig{
		Service: r.service,
		Logger:  r.log,
	}))

	r.Handle("PATCH /v1/{id}", v1.NewUpdateHandler(&v1.UpdateHandlerConfig{
		Service: r.service,
		Logger:  r.log,
	}))

	r.Handle("DELETE /v1/{id}", v1.NewDeleteHandler(&v1.DeleteHandlerConfig{
		Service: r.service,
		Logger:  r.log,
	}))
}

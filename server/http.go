package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/service"
	"gorm.io/gorm"
)

type HTTPServer struct {
	echo *echo.Echo

	//	Port
	port string

	//	Database Dialector
	dialector gorm.Dialector

	//	Logger
	logger *slog.Logger
}

type NewHTTPServerConfig struct {

	//	Port
	Port string

	//	Database Dialector
	Dialector gorm.Dialector

	//	Logger
	Logger *slog.Logger
}

func NewHTTPServer(config *NewHTTPServerConfig) *HTTPServer {

	server := HTTPServer{
		echo:      echo.New(),
		port:      config.Port,
		dialector: config.Dialector,
		logger:    config.Logger,
	}

	if server.logger != nil {

		//	Setup the logger.
		server.logger = slog.Default()
	}

	//	Add default middlewares.
	server.echo.Use(middleware.Recover())

	server.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				server.logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				server.logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("error", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	return &server
}

// Serve starts the HTTP server.
func (s *HTTPServer) Serve() {
	s.echo.Logger.Fatal(s.echo.Start(fmt.Sprintf(":%s", s.port)))
}

func (s *HTTPServer) InitDefaultRoutes(base string) {

	group := s.echo.Group(base)

	//	Setup a health check endpoint.
	group.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	//	API v1 routes.
	v1 := group.Group("/v1")

	v1.POST("", s.create)
	v1.GET("", s.list)
	v1.GET("/:id", s.get)
	v1.PATCH("/:id", s.update)
	v1.DELETE("/:id", s.delete)
}

// Open database connection and initialize the service.
func (s *HTTPServer) getService() (service.Service, error) {

	//	Prepare a database connection.
	database, err := gorm.Open(s.dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//	Get the svc.
	svc := service.NewService(&service.Config{
		DB: db.NewDB(&db.Config{
			DB: database,
		}),
	})

	return svc, nil
}

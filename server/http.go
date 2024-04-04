package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/service"
	"gorm.io/gorm"

	slogGorm "github.com/orandin/slog-gorm"
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

		//	Setup a new logger.
		server.logger = slog.Default()
	}

	//
	//	Add default middlewares.
	//

	//	Recover middleware recovers from panics anywhere in the chain, prints stack trace and handovers the control to the centralized HTTPErrorHandler.
	server.echo.Use(middleware.Recover())

	//	Request ID middleware generates a unique ID for every request.
	server.echo.Use(middleware.RequestID())

	//	RequestLogger middleware allows developer fully to customize what is logged and how it is logged and is more suitable for usage with 3rd party (structured logging) libraries.
	server.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				server.logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "HTTP",
					slog.String("id", c.Response().Header().Get(echo.HeaderXRequestID)),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				server.logger.LogAttrs(c.Request().Context(), slog.LevelError, "HTTP_ERROR",
					slog.String("id", c.Response().Header().Get(echo.HeaderXRequestID)),
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
		return c.JSON(http.StatusOK, &Response{
			Message: "Service is up and running.",
			Data: map[string]interface{}{
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			},
		})
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

	gormLogger := slogGorm.New(
		slogGorm.WithHandler(s.logger.Handler()), // since v1.3.0
		slogGorm.WithTraceAll(),                  // trace all messages
	)

	//	Prepare a database connection.
	database, err := gorm.Open(s.dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	//	Setup the database service.
	db := db.NewDB(&db.Config{
		DB: database,
	})

	//	Setup the business service.
	svc := service.NewService(&service.Config{
		Logger: s.logger,
		DB:     db,
	})

	return svc, nil
}

//
//	Handlers
//

// Create Handler.
func (s *HTTPServer) create(c echo.Context) error {

	//	Unmarshal the incoming payload.
	var payload CreateOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Initialize a default context.
	ctx := c.Request().Context()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todo, err := svc.Create(ctx, &service.CreateOptions{
		Title: payload.Title,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to create the todo.",
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo created successfully.",
		Data:    todo,
	})
}

// Get Handler.
func (s *HTTPServer) get(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Initialize a default context.
	ctx := c.Request().Context()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todo, err := svc.Get(ctx, uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo fetched successfully.",
		Data:    todo,
	})
}

// List Handler.
func (s *HTTPServer) list(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload ListOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Initialize a default context.
	ctx := c.Request().Context()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todos, err := svc.List(ctx, &service.ListOptions{
		Skip:           payload.Skip,
		Limit:          payload.Limit,
		Title:          payload.Title,
		OrderBy:        payload.OrderBy,
		OrderDirection: payload.OrderDirection,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to list the todos.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todos fetched successfully.",
		Data:    todos,
	})
}

// Update Handler.
func (s *HTTPServer) update(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Unmarshal the incoming payload.
	var payload UpdateOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Initialize a default context.
	ctx := c.Request().Context()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todo, err := svc.Update(ctx, uuid, &service.UpdateOptions{
		Title: payload.Title,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo updated successfully.",
		Data:    todo,
	})
}

// Delete Handler.
func (s *HTTPServer) delete(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Initialize a default context.
	ctx := c.Request().Context()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	err = svc.Delete(ctx, uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo deleted successfully.",
	})
}

package router

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/service"
	"gorm.io/gorm"

	slogGorm "github.com/orandin/slog-gorm"
)

type HTTPRouter struct {
	*http.ServeMux

	// Gorm database dialector to use.
	// Example: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	//
	// This field is mandatory.
	dialector gorm.Dialector

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

	// Gorm database dialector to use.
	// Example: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	//
	// This field is mandatory.
	Dialector gorm.Dialector

	// Logger is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	Logger *slog.Logger
}

// NewHTTPRouter creates a new instance of `HTTPRouter`.
func NewHTTPRouter(config *HTTPRouterConfig) *HTTPRouter {

	router := HTTPRouter{
		ServeMux:  http.NewServeMux(),
		dialector: config.Dialector,
		log:       config.Logger,
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
	router.HandleFunc("POST /", handle(router.Create))

	return &router
}

// OpenDB opens a database connection and returns the db/respository service.
func (r *HTTPRouter) OpenDB() (db.DB, error) {

	//	Setup the gorm logger.
	handler := r.log.With("layer", "database").Handler()
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(handler), // since v1.3.0
		slogGorm.WithTraceAll(),       // trace all messages
	)

	//	Prepare a database connection.
	database, err := gorm.Open(r.dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	//	Setup the database service.
	return db.NewDB(&db.Config{
		DB: database,
	}), nil
}

//
// Functions which will handle incoming requests.
//

// Create handler create a new record.
func (r *HTTPRouter) Create(req *http.Request) error {

	// Decode the request body.
	body, err := decode[CreateOptions](req)
	if err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Err:     err,
		}
	}

	// Open the database connection.
	db, err := r.OpenDB()
	if err != nil {
		return &Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to open the database connection.",
			Err:     err,
		}
	}

	// Get the appropriate business service.
	svc := service.NewService(&service.Config{
		DB:     db,
		Logger: r.log,
	})

	// Call the service method that performs the required operation.
	record, err := svc.Create(req.Context(), &service.CreateOptions{
		Title: body.Title,
	})
	if err != nil {
		return &Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create the record.",
			Err:     err,
		}
	}

	return &Response{
		Status:  http.StatusCreated,
		Message: "The record was created successfully.",
		Data:    record,
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
				if err := write(w, response.Status, response); err != nil {
					log.Println("failed to write response:", err)
				}
				return
			}
			if err := write(w, http.StatusInternalServerError, &Response{
				Message: "Your broke something on our server :(",
				Err:     err,
			}); err != nil {
				log.Println("failed to write response:", err)
			}
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

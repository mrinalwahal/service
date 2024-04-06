package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/service"
	"gorm.io/gorm"
)

// CreateOptions represents the options for creating a todo.
type CreateOptions struct {

	//	Title of the todo.
	Title string `json:"title" validate:"required"`
}

// Create handler create a new record.
type CreateHandler struct {

	// Gorm database dialector to use.
	// Example: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	//
	// This field is mandatory.
	dialector gorm.Dialector

	// Database connection based on dialector.
	// We would ideally open this at the time of serving the request and keep it open for all base functions to use it.
	db db.DB

	// Options contains the payload received in the incoming request.
	// This is useful in passing the request payload to the service layer.
	// For example, it contains the request body in case of a POST request. Or the query parameters in case of a GET request.
	options *CreateOptions

	// log is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	log *slog.Logger
}

type CreateHandlerConfig struct {

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

// NewCreateHandler creates a new instance of `CreateHandler`.
func NewCreateHandler(config *CreateHandlerConfig) *CreateHandler {
	handler := CreateHandler{
		log:       config.Logger,
		dialector: config.Dialector,
	}

	// Set the default logger if not provided.
	if handler.log == nil {
		handler.log = slog.Default()
	}

	handler.log = handler.log.With("handler", "create")

	return &handler
}

// ServeHTTP handles the incoming HTTP request.
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Decode the request options.
	options, err := decode[CreateOptions](r)
	if err != nil {
		write(w, http.StatusInternalServerError, &response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Err:     err,
		})
		return
	}
	h.options = &options

	// Load the context.
	ctx := r.Context()

	// Open the database connection.
	db, err := db.NewDB(&db.Config{
		Dialector: h.dialector,
		Logger:    h.log,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, &response{
			Message: "Failed to connect to the database.",
		})
		return
	}
	h.db = db

	// Validate the request.
	if err := h.validate(ctx); err != nil {
		handleErr(w, err)
		return
	}

	// Call the function.
	if err := h.function(ctx); err != nil {
		handleErr(w, err)
	}
}

// validate function ascertains that the requester is authorized to perform this request.
// This is where the "API rule/condition" logic is applied.
func (h *CreateHandler) validate(ctx context.Context) error {
	return nil
}

// function applies the fundamental business logic to complete required operation.
func (h *CreateHandler) function(ctx context.Context) error {

	// Get the appropriate business service.
	svc := service.NewService(&service.Config{
		DB:     h.db,
		Logger: h.log,
	})

	// Call the service method that performs the required operation.
	record, err := svc.Create(ctx, &service.CreateOptions{
		Title: h.options.Title,
	})
	if err != nil {
		return &response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create the record.",
			Err:     err,
		}
	}

	return &response{
		Status:  http.StatusCreated,
		Message: "The record was created successfully.",
		Data:    record,
	}
}

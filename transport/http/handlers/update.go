package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/service"
)

// UpdateOptions represents the options for updating a record.
type UpdateOptions struct {

	//	Title of the record.
	Title string `json:"title" validate:"required"`
}

// Update handler update a new record.
type UpdateHandler struct {

	// Database layer.
	// The connection should already be open.
	//
	// This field is mandatory.
	db db.DB

	// The UUID of the record to update.
	id uuid.UUID

	// Options contains the payload received in the incoming request.
	// This is useful in passing the request payload to the service layer.
	// For example, it contains the request body in case of a POST request. Or the query parameters in case of a GET request.
	options *UpdateOptions

	// log is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	log *slog.Logger
}

type UpdateHandlerConfig struct {

	// Database layer.
	// The connection should already be open.
	//
	// This field is mandatory.
	DB db.DB

	// Logger is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	Logger *slog.Logger
}

// NewUpdateHandler updates a new instance of `UpdateHandler`.
func NewUpdateHandler(config *UpdateHandlerConfig) *UpdateHandler {
	handler := UpdateHandler{
		db:  config.DB,
		log: config.Logger,
	}

	// Set the default logger if not provided.
	if handler.log == nil {
		handler.log = slog.Default()
	}
	handler.log = handler.log.With("handler", "update")

	return &handler
}

// ServeHTTP handles the incoming HTTP request.
func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Decode the request options.
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		write(w, http.StatusBadRequest, &response{
			Message: "Invalid ID.",
		})
		return
	}
	h.id = id

	options, err := decode[UpdateOptions](r)
	if err != nil {
		write(w, http.StatusBadRequest, &response{
			Message: "Invalid request options.",
			Err:     err,
		})
		return
	}
	h.options = &options

	// Load the context.
	ctx := r.Context()

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
func (h *UpdateHandler) validate(ctx context.Context) error {
	return nil
}

// function applies the fundamental business logic to complete required operation.
func (h *UpdateHandler) function(ctx context.Context) error {

	// Get the appropriate business service.
	svc := service.NewService(&service.Config{
		DB:     h.db,
		Logger: h.log,
	})

	// Call the service method that performs the required operation.
	record, err := svc.Update(ctx, h.id, &service.UpdateOptions{
		Title: h.options.Title,
	})
	if err != nil {
		return &response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update the record.",
			Err:     err,
		}
	}

	return &response{
		Status:  http.StatusOK,
		Message: "The record was updated successfully.",
		Data:    record,
	}
}

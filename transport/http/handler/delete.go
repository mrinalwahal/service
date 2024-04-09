package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/service"
)

// Delete handler deletes the record.
type DeleteHandler struct {

	// Service layer.
	//
	// This field is mandatory.
	service service.Service

	// The UUID of the record to delete.
	id uuid.UUID

	// log is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	log *slog.Logger
}

type DeleteHandlerConfig struct {

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

// NewDeleteHandler deletes a new instance of `DeleteHandler`.
func NewDeleteHandler(config *DeleteHandlerConfig) Handler {
	handler := DeleteHandler{
		service: config.Service,
		log:     config.Logger,
	}

	// Set the default logger if not provided.
	if handler.log == nil {
		handler.log = slog.Default()
	}
	handler.log = handler.log.With("handler", "delete")

	return &handler
}

// ServeHTTP handles the incoming HTTP request.
func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Decode the request options.
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Invalid ID.",
		})
		return
	}
	h.id = id

	// Load the context.
	ctx := r.Context()

	// Validate the request.
	if err := h.validate(ctx, id); err != nil {
		handleErr(w, err)
		return
	}

	// Call the function.
	if err := h.process(ctx, id); err != nil {
		handleErr(w, err)
	}
}

// Validate function ascertains that the requester is authorized to perform this request.
// This is where the "API rule/condition" logic is applied.
func (h *DeleteHandler) validate(ctx context.Context, ID uuid.UUID) error {
	return nil
}

// Process applies the fundamental business logic to complete required operation.
func (h *DeleteHandler) process(ctx context.Context, ID uuid.UUID) error {

	// Call the service method that performs the required operation.
	if err := h.service.Delete(ctx, ID); err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to delete the record.",
			Err:     err,
		}
	}

	return &Response{
		Status:  http.StatusOK,
		Message: "The record was deleted successfully.",
	}
}

package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/service"
)

// UpdateOptions represents the options for updating a record.
type UpdateOptions struct {

	//	Title of the record.
	Title string `json:"title" validate:"required"`
}

// Update handler update a new record.
type UpdateHandler struct {

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

type UpdateHandlerConfig struct {

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

// NewUpdateHandler updates a new instance of `UpdateHandler`.
func NewUpdateHandler(config *UpdateHandlerConfig) Handler {
	handler := UpdateHandler{
		service: config.Service,
		log:     config.Logger,
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
		write(w, http.StatusBadRequest, &Response{
			Message: "Invalid ID.",
		})
		return
	}

	options, err := decode[UpdateOptions](r)
	if err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Invalid request options.",
			Err:     err,
		})
		return
	}

	// Load the context.
	ctx := r.Context()

	// Validate the request.
	if err := h.validate(ctx, id, &options); err != nil {
		handleErr(w, err)
		return
	}

	// Call the function.
	if err := h.process(ctx, id, &options); err != nil {
		handleErr(w, err)
	}
}

// validate function ascertains that the requester is authorized to perform this request.
// This is where the "API rule/condition" logic is applied.
func (h *UpdateHandler) validate(ctx context.Context, ID uuid.UUID, options *UpdateOptions) error {
	return nil
}

// process applies the fundamental business logic to complete required operation.
func (h *UpdateHandler) process(ctx context.Context, ID uuid.UUID, options *UpdateOptions) error {

	// Call the service method that performs the required operation.
	record, err := h.service.Update(ctx, ID, &service.UpdateOptions{
		Title: options.Title,
	})
	if err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to update the record.",
			Err:     err,
		}
	}

	return &Response{
		Status:  http.StatusOK,
		Message: "The record was updated successfully.",
		Data:    record,
	}
}

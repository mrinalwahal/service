package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/service"
)

// Get handler gets the record.
type GetHandler struct {

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

type GetHandlerConfig struct {

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

// NewGetHandler gets a new instance of `GetHandler`.
func NewGetHandler(config *GetHandlerConfig) Handler {
	handler := GetHandler{
		service: config.Service,
		log:     config.Logger,
	}

	// Set the default logger if not provided.
	if handler.log == nil {
		handler.log = slog.Default()
	}
	handler.log = handler.log.With("handler", "get")

	return &handler
}

// ServeHTTP handles the incoming HTTP request.
func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Decode the request options.
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Invalid ID.",
		})
		return
	}

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

// validate function ascertains that the requester is authorized to perform this request.
// This is where the "API rule/condition" logic is applied.
func (h *GetHandler) validate(ctx context.Context, ID uuid.UUID) error {
	return nil
}

// process applies the fundamental business logic to complete required operation.
func (h *GetHandler) process(ctx context.Context, ID uuid.UUID) error {

	// Call the service method that performs the required operation.
	record, err := h.service.Get(ctx, ID)
	if err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to get the record.",
			Err:     err,
		}
	}

	return &Response{
		Status:  http.StatusOK,
		Message: "The record was retrieved successfully.",
		Data:    record,
	}
}

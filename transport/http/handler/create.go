package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mrinalwahal/service/service"
)

// CreateOptions represents the options for creating a record.
type CreateOptions struct {

	//	Title of the record.
	Title string `json:"title"`
}

// Create handler create a new record.
type CreateHandler struct {

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

type CreateHandlerConfig struct {

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

// NewCreateHandler creates a new instance of `CreateHandler`.
func NewCreateHandler(config *CreateHandlerConfig) Handler {
	handler := CreateHandler{
		service: config.Service,
		log:     config.Logger,
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
		write(w, http.StatusBadRequest, &response{
			Message: "Invalid request options.",
			Err:     err,
		})
		return
	}

	// Load the context.
	ctx := r.Context()

	// Validate the request.
	if err := h.validate(ctx, &options); err != nil {
		handleErr(w, err)
		return
	}

	// Call the function.
	if err := h.process(ctx, &options); err != nil {
		handleErr(w, err)
	}
}

// validate function ascertains that the requester is authorized to perform this request.
// This is where the "API rule/condition" logic is applied.
func (h *CreateHandler) validate(ctx context.Context, options *CreateOptions) error {
	return nil
}

// process applies the fundamental business logic to complete required operation.
func (h *CreateHandler) process(ctx context.Context, options *CreateOptions) error {

	// Call the service method that performs the required operation.
	record, err := h.service.Create(ctx, &service.CreateOptions{
		Title: options.Title,
	})
	if err != nil {
		return &response{
			Status:  http.StatusBadRequest,
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

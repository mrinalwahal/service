package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/pkg/middleware"
	"github.com/mrinalwahal/service/service"
)

// CreateOptions represents the options for creating a record.
type CreateOptions struct {

	//	Title of the record.
	Title string `json:"title"`

	// UserID extracted from the request context.
	UserID uuid.UUID `json:"-"`
}

// Validate the options.
func (o *CreateOptions) Validate() error {
	if o.Title == "" {
		return &Response{
			Status:  http.StatusBadRequest,
			Message: "Title is required.",
			Err:     ErrInvalidRequestOptions,
		}
	}
	return nil
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
		write(w, http.StatusBadRequest, &Response{
			Message: "Invalid request options.",
			Err:     err,
		})
		return
	}

	// Load the context.
	ctx := r.Context()

	// Load the claims from request context to pass them in the service method.
	userID, ok := ctx.Value(middleware.UserID).(string)
	if !ok {
		handleErr(w, &Response{
			Status:  http.StatusUnauthorized,
			Message: "User ID not found in the request context.",
			Err:     ErrInvalidRequestOptions,
		})
		return
	}
	options.UserID, err = uuid.Parse(userID)
	if err != nil {
		handleErr(w, &Response{
			Status:  http.StatusBadRequest,
			Message: "Failed to parse the user ID.",
			Err:     err,
		})
		return
	}

	// Validate the request options.
	if err := options.Validate(); err != nil {
		handleErr(w, err)
		return
	}

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
		Title:  options.Title,
		UserID: options.UserID,
	})
	if err != nil {
		return &Response{
			Status:  http.StatusBadRequest,
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

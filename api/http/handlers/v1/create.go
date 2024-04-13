package v1

import (
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

	// ID of the user who is creating the record.
	UserID uuid.UUID `json:"-"`
}

// validate the options.
func (o *CreateOptions) validate() error {
	if o.Title == "" {
		return ErrInvalidRequestOptions
	}
	return nil
}

// preset presets options from claims in the request context.
func (o *CreateOptions) preset(r *http.Request) error {
	claims, exists := r.Context().Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if !exists {
		return ErrInvalidJWTClaims
	}

	o.UserID = claims.XUserID
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
	h.log.DebugContext(r.Context(), "handling request")

	// Decode the request options.
	options, err := decode[CreateOptions](r)
	if err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Invalid request options.",
			Err:     err,
		})
		return
	}

	// Validate the request options.
	if err := options.validate(); err != nil {
		write(w, http.StatusBadRequest, Response{
			Message: "Failed validate request options.",
			Err:     ErrInvalidRequestOptions,
		})
		return
	}

	// Present options from the request.
	if err := options.preset(r); err != nil {
		write(w, http.StatusBadRequest, Response{
			Message: "Failed to apply presets.",
			Err:     err,
		})
		return
	}

	// Call the service method that performs the required operation.
	record, err := h.service.Create(r.Context(), &service.CreateOptions{
		Title:  options.Title,
		UserID: options.UserID,
	})
	if err != nil {
		write(w, http.StatusBadRequest, Response{
			Message: "Failed to create the record.",
			Err:     err,
		})
		return
	}

	write(w, http.StatusCreated, Response{
		Message: "The record was created successfully.",
		Data:    record,
	})
}

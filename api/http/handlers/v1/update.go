package v1

import (
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

	record, err := h.service.Update(r.Context(), id, &service.UpdateOptions{
		Title: options.Title,
	})
	if err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Failed to update the record.",
			Err:     err,
		})
		return
	}

	write(w, http.StatusOK, &Response{
		Message: "The record was updated successfully.",
		Data:    record,
	})
	return
}

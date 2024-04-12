package v1

import (
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
			Err:     err,
		})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Failed to delete the record.",
			Err:     err,
		})
		return
	}

	write(w, http.StatusOK, &Response{
		Message: "The record was deleted successfully.",
	})
}

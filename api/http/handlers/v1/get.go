package v1

import (
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

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Invalid ID.",
		})
		return
	}

	record, err := h.service.Get(r.Context(), id)
	if err != nil {
		write(w, http.StatusBadRequest, &Response{
			Message: "Failed to get the record.",
			Err:     err,
		})
		return
	}

	write(w, http.StatusOK, &Response{
		Message: "The record was retrieved successfully.",
		Data:    record,
	})
}

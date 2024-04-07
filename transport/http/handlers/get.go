package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/service"
)

// Get handler gets the record.
type GetHandler struct {

	// Database layer.
	// The connection should already be open.
	//
	// This field is mandatory.
	db db.DB

	// The UUID of the record to get.
	id uuid.UUID

	// log is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	log *slog.Logger
}

type GetHandlerConfig struct {

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

// NewGetHandler gets a new instance of `GetHandler`.
func NewGetHandler(config *GetHandlerConfig) *GetHandler {
	handler := GetHandler{
		db:  config.DB,
		log: config.Logger,
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
		write(w, http.StatusBadRequest, &response{
			Message: "Invalid ID.",
		})
		return
	}
	h.id = id

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
func (h *GetHandler) validate(ctx context.Context) error {
	return nil
}

// function applies the fundamental business logic to complete required operation.
func (h *GetHandler) function(ctx context.Context) error {

	// Get the appropriate business service.
	svc := service.NewService(&service.Config{
		DB:     h.db,
		Logger: h.log,
	})

	// Call the service method that performs the required operation.
	record, err := svc.Get(ctx, h.id)
	if err != nil {
		return &response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get the record.",
			Err:     err,
		}
	}

	return &response{
		Status:  http.StatusOK,
		Message: "The record was retrieved successfully.",
		Data:    record,
	}
}

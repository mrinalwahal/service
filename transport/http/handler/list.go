package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dyninc/qstring"
	"github.com/mrinalwahal/service/service"
)

// ListOptions represents the options for listing records.
type ListOptions struct {

	//	Number of records to skip.
	Skip int `query:"skip" validate:"gte=0"`

	//	Number of records to return.
	Limit int `query:"limit" validate:"gte=0,lte=100"`

	//	Order by field.
	OrderBy string `query:"orderBy" validate:"oneof=created_at updated_at title"`

	//	Order by direction.
	OrderDirection string `query:"orderDirection" validate:"oneof=asc desc"`

	//	Title of the record.
	Title string `query:"name"`
}

// List handler lists the records.
type ListHandler struct {

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

type ListHandlerConfig struct {

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

// NewListHandler lists a new instance of `ListHandler`.
func NewListHandler(config *ListHandlerConfig) Handler {
	handler := ListHandler{
		service: config.Service,
		log:     config.Logger,
	}

	// Set the default logger if not provided.
	if handler.log == nil {
		handler.log = slog.Default()
	}
	handler.log = handler.log.With("handler", "list")

	return &handler
}

// ServeHTTP handles the incoming HTTP request.
func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Decode the request options.
	var options ListOptions
	err := qstring.Unmarshal(r.URL.Query(), &options)
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

// Validate function ascertains that the requester is authorized to perform this request.
// This is where the "API rule/condition" logic is applied.
func (h *ListHandler) validate(ctx context.Context, options *ListOptions) error {
	return nil
}

// process applies the fundamental business logic to complete required operation.
func (h *ListHandler) process(ctx context.Context, options *ListOptions) error {

	// Call the service method that performs the required operation.
	records, err := h.service.List(ctx, &service.ListOptions{
		Title: options.Title,
	})
	if err != nil {
		return &response{
			Status:  http.StatusBadRequest,
			Message: "Failed to list the records.",
			Err:     err,
		}
	}

	return &response{
		Status:  http.StatusOK,
		Message: "The records were retrieved successfully.",
		Data:    records,
	}
}

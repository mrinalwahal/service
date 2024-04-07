package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dyninc/qstring"
	"github.com/mrinalwahal/service/db"
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

	// Database layer.
	// The connection should already be open.
	//
	// This field is mandatory.
	db db.DB

	// Options contains the payload received in the incoming request.
	// This is useful in passing the request payload to the service layer.
	// For example, it contains the request body in case of a POST request. Or the query parameters in case of a GET request.
	options *ListOptions

	// log is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	log *slog.Logger
}

type ListHandlerConfig struct {

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

// NewListHandler lists a new instance of `ListHandler`.
func NewListHandler(config *ListHandlerConfig) Handler {
	handler := ListHandler{
		db:  config.DB,
		log: config.Logger,
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
	h.options = &options

	// Load the context.
	ctx := r.Context()

	// Validate the request.
	if err := h.Validate(ctx); err != nil {
		handleErr(w, err)
		return
	}

	// Call the function.
	if err := h.Process(ctx); err != nil {
		handleErr(w, err)
	}
}

// Validate function ascertains that the requester is authorized to perform this request.
// This is where the "API rule/condition" logic is applied.
func (h *ListHandler) Validate(ctx context.Context) error {
	return nil
}

// Process applies the fundamental business logic to complete required operation.
func (h *ListHandler) Process(ctx context.Context) error {

	// Get the appropriate business service.
	svc := service.NewService(&service.Config{
		DB:     h.db,
		Logger: h.log,
	})

	// Call the service method that performs the required operation.
	records, err := svc.List(ctx, &service.ListOptions{
		Title: h.options.Title,
	})
	if err != nil {
		return &response{
			Status:  http.StatusInternalServerError,
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

package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dyninc/qstring"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/service"
	"gorm.io/gorm"
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

	// Gorm database dialector to use.
	// Example: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	//
	// This field is mandatory.
	dialector gorm.Dialector

	// Database connection based on dialector.
	// We would ideally open this at the time of serving the request and keep it open for all base functions to use it.
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

	// Gorm database dialector to use.
	// Example: postgres.Open("host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Kolkata")
	//
	// This field is mandatory.
	Dialector gorm.Dialector

	// Logger is the `log/slog` instance that will be used to log messages.
	// Default: `slog.DefaultLogger`
	//
	// This field is optional.
	Logger *slog.Logger
}

// NewListHandler lists a new instance of `ListHandler`.
func NewListHandler(config *ListHandlerConfig) *ListHandler {
	handler := ListHandler{
		log:       config.Logger,
		dialector: config.Dialector,
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

	// Open the database connection.
	db, err := db.NewDB(&db.Config{
		Dialector: h.dialector,
		Logger:    h.log,
	})
	if err != nil {
		write(w, http.StatusInternalServerError, &response{
			Message: "Failed to connect to the database.",
		})
		return
	}
	h.db = db

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
func (h *ListHandler) validate(ctx context.Context) error {
	return nil
}

// function applies the fundamental business logic to complete required operation.
func (h *ListHandler) function(ctx context.Context) error {

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

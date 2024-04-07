package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/db"
)

type Service interface {
	Create(context.Context, *CreateOptions) (*db.Record, error)
	List(context.Context, *ListOptions) ([]*db.Record, error)
	Get(context.Context, uuid.UUID) (*db.Record, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*db.Record, error)
	Delete(context.Context, uuid.UUID) error
}

type Config struct {

	//	Database layer service.
	DB db.DB

	//	Logger.
	Logger *slog.Logger
}

// Initializes and gets the service with the supplied database connection.
func NewService(config *Config) Service {
	svc := service{
		db:     config.DB,
		logger: config.Logger,
	}

	if svc.logger == nil {
		svc.logger = slog.Default()
	}

	svc.logger = svc.logger.With("layer", "service")

	return &svc
}

type service struct {

	//	Database layer service.
	db db.DB

	//	Logger.
	logger *slog.Logger
}

func (s *service) Create(ctx context.Context, options *CreateOptions) (*db.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "creating a new record",
		slog.String("function", "create"),
	)
	return s.db.Create(ctx, &db.CreateOptions{
		Title: options.Title,
	})
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*db.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "listing all records",
		slog.String("function", "list"),
	)
	return s.db.List(ctx, &db.ListOptions{
		Title:          options.Title,
		Skip:           options.Skip,
		Limit:          options.Limit,
		OrderBy:        options.OrderBy,
		OrderDirection: options.OrderDirection,
	})
}

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*db.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "retrieving a record",
		slog.String("function", "get"),
	)
	return s.db.Get(ctx, ID)
}

func (s *service) Update(ctx context.Context, ID uuid.UUID, options *UpdateOptions) (*db.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "updating a record",
		slog.String("function", "update"),
	)
	return s.db.Update(ctx, ID, &db.UpdateOptions{
		Title: options.Title,
	})
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "deleting a record",
		slog.String("function", "delete"),
	)
	return s.db.Delete(ctx, ID)
}

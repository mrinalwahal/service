package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/db"
)

type Service interface {
	Create(context.Context, *CreateOptions) (*db.Todo, error)
	List(context.Context, *ListOptions) ([]*db.Todo, error)
	Get(context.Context, uuid.UUID) (*db.Todo, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*db.Todo, error)
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

	if svc.logger != nil {
		svc.logger = svc.logger.With("layer", "service")
	} else {
		svc.logger = slog.Default()
	}

	return &svc
}

type service struct {

	//	Database layer service.
	db db.DB

	//	Logger.
	logger *slog.Logger
}

func (s *service) Create(ctx context.Context, options *CreateOptions) (*db.Todo, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "creating a new todo",
		slog.String("function", "create"),
	)
	return s.db.Create(ctx, &db.CreateOptions{
		Title: options.Title,
	})
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*db.Todo, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "listing all todos",
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

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*db.Todo, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "retrieving a todo",
		slog.String("function", "get"),
	)
	return s.db.Get(ctx, ID)
}

func (s *service) Update(ctx context.Context, ID uuid.UUID, options *UpdateOptions) (*db.Todo, error) {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "updating a todo",
		slog.String("function", "update"),
	)
	return s.db.Update(ctx, ID, &db.UpdateOptions{
		Title: options.Title,
	})
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	s.logger.LogAttrs(ctx, slog.LevelInfo, "deleting a todo",
		slog.String("function", "delete"),
	)
	return s.db.Delete(ctx, ID)
}

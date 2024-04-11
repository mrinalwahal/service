//go:generate mockgen -destination=mock.go -source=service.go -package=service
package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/commons"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/model"
)

type Service interface {
	Create(context.Context, *CreateOptions, *commons.Requester) (*model.Record, error)
	List(context.Context, *ListOptions, *commons.Requester) ([]*model.Record, error)
	Get(context.Context, uuid.UUID, *commons.Requester) (*model.Record, error)
	Update(context.Context, uuid.UUID, *UpdateOptions, *commons.Requester) (*model.Record, error)
	Delete(context.Context, uuid.UUID, *commons.Requester) error
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

func (s *service) Create(ctx context.Context, options *CreateOptions, requester *commons.Requester) (*model.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "creating a new record",
		slog.String("function", "create"),
	)
	if options == nil {
		return nil, ErrOptionsNotFound
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	return s.db.Create(ctx, &db.CreateOptions{
		Title:  options.Title,
		UserID: options.UserID,
	}, requester)
}

func (s *service) List(ctx context.Context, options *ListOptions, requester *commons.Requester) ([]*model.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "listing all records",
		slog.String("function", "list"),
	)
	if options == nil {
		return nil, ErrOptionsNotFound
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	return s.db.List(ctx, &db.ListOptions{
		Title:          options.Title,
		Skip:           options.Skip,
		Limit:          options.Limit,
		OrderBy:        options.OrderBy,
		OrderDirection: options.OrderDirection,
	}, &db.Requester{
		ID: requester.ID,
	})
}

func (s *service) Get(ctx context.Context, ID uuid.UUID, requester *commons.Requester) (*model.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "retrieving a record",
		slog.String("function", "get"),
	)
	if ID == uuid.Nil {
		return nil, ErrOptionsNotFound
	}
	return s.db.Get(ctx, ID, &db.Requester{
		ID: requester.ID,
	})
}

func (s *service) Update(ctx context.Context, ID uuid.UUID, options *UpdateOptions, requester *commons.Requester) (*model.Record, error) {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "updating a record",
		slog.String("function", "update"),
	)
	if ID == uuid.Nil {
		return nil, ErrOptionsNotFound
	}
	if options == nil {
		return nil, ErrOptionsNotFound
	}
	if err := options.validate(); err != nil {
		return nil, err
	}
	return s.db.Update(ctx, ID, &db.UpdateOptions{
		Title: options.Title,
	}, &db.Requester{
		ID: requester.ID,
	})
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID, requester *commons.Requester) error {
	s.logger.LogAttrs(ctx, slog.LevelDebug, "deleting a record",
		slog.String("function", "delete"),
	)
	if ID == uuid.Nil {
		return ErrOptionsNotFound
	}
	return s.db.Delete(ctx, ID, &db.Requester{
		ID: requester.ID,
	})
}

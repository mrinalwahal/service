package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/db"
)

type Service interface {
	Create(context.Context, *CreateOptions) (*db.Todo, error)
	Get(context.Context, uuid.UUID) (*db.Todo, error)
	List(context.Context, *ListOptions) ([]*db.Todo, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*db.Todo, error)
	Delete(context.Context, uuid.UUID) error
}

type Config struct {

	//	Database layer service.
	DB db.DB
}

// Initializes and gets the service with the supplied database connection.
func NewService(config *Config) Service {
	return &service{
		db: config.DB,
	}
}

type service struct {

	//	Database layer service.
	db db.DB
}

func (s *service) Create(ctx context.Context, options *CreateOptions) (*db.Todo, error) {
	return s.db.Create(ctx, &db.CreateOptions{
		Title: options.Title,
	})
}

func (s *service) Get(ctx context.Context, ID uuid.UUID) (*db.Todo, error) {
	return s.db.Get(ctx, ID)
}

func (s *service) List(ctx context.Context, options *ListOptions) ([]*db.Todo, error) {
	return s.db.List(ctx, &db.ListOptions{
		Title:          options.Title,
		Skip:           options.Skip,
		Limit:          options.Limit,
		OrderBy:        options.OrderBy,
		OrderDirection: options.OrderDirection,
	})
}

func (s *service) Update(ctx context.Context, ID uuid.UUID, options *UpdateOptions) (*db.Todo, error) {
	return s.db.Update(ctx, ID, &db.UpdateOptions{
		Title: options.Title,
	})
}

func (s *service) Delete(ctx context.Context, ID uuid.UUID) error {
	return s.db.Delete(ctx, ID)
}

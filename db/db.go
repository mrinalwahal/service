//go:generate mockgen -destination=db_mock.go -source=db.go -package=db
package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
)

// DB interface declares the signature of the database layer.
type DB interface {
	Create(context.Context, *CreateOptions) (*model.Record, error)
	List(context.Context, *ListOptions) ([]*model.Record, error)
	Get(context.Context, uuid.UUID) (*model.Record, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*model.Record, error)
	Delete(context.Context, uuid.UUID) error
}

//go:generate mockgen -destination=mock.go -source=db.go -package=db
package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
	"gorm.io/gorm"
)

// DB interface declares the signature of the database layer.
type DB interface {
	Create(context.Context, *CreateOptions) (*model.Record, error)
	Get(context.Context, uuid.UUID) (*model.Record, error)
	List(context.Context, *ListOptions) ([]*model.Record, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*model.Record, error)
	Delete(context.Context, uuid.UUID) error
}

type Config struct {

	// Database connection.
	// The connection should already be open.
	//
	// This field is mandatory.
	DB *gorm.DB
}

func NewDB(config *Config) DB {
	db := database{
		conn: config.DB,
	}

	return &db
}

type database struct {

	//	Database Connection
	conn *gorm.DB
}

// Create operation creates a new record in the database.
func (db *database) Create(ctx context.Context, options *CreateOptions) (*model.Record, error) {
	txn := db.conn.WithContext(ctx)
	// if err := options.validate(); err != nil {
	// 	return nil, err
	// }

	// Prepare the payload we have to send to the database transaction.
	var payload model.Record
	payload.Title = options.Title
	payload.UserID = options.UserID

	// Execute the transaction.
	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *database) List(ctx context.Context, options *ListOptions) ([]*model.Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload []*model.Record

	query := txn
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}
	if options.Skip > 0 {
		query = query.Offset(options.Skip)
	}
	if options.OrderBy != "" {
		query = query.Order(options.OrderBy + " " + options.OrderDirection)
	}
	if options.Title != "" {
		query = query.Where(&model.Record{
			Title: options.Title,
		})
	}
	if options.UserID != uuid.Nil {
		query = query.Where(&model.Record{
			UserID: options.UserID,
		})
	}

	if result := query.Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

func (db *database) Get(ctx context.Context, ID uuid.UUID) (*model.Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload model.Record
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *database) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*model.Record, error) {
	txn := db.conn.WithContext(ctx)
	if err := options.validate(); err != nil {
		return nil, err
	}

	var payload model.Record
	payload.ID = id
	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return db.Get(ctx, id)
}

func (db *database) Delete(ctx context.Context, ID uuid.UUID) error {
	txn := db.conn.WithContext(ctx)

	var payload model.Record
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}

package db

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DB interface contains the signature of the database layer.
type DB interface {
	Create(context.Context, *CreateOptions) (*Record, error)
	Get(context.Context, uuid.UUID) (*Record, error)
	List(context.Context, *ListOptions) ([]*Record, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*Record, error)
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

func (db *database) Create(ctx context.Context, options *CreateOptions) (*Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.Title = options.Title

	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *database) List(ctx context.Context, options *ListOptions) ([]*Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload []*Record

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

	//	Add conditions to the query.
	where := Record{
		Title: options.Title,
	}

	if result := query.Where(&where).Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

func (db *database) Get(ctx context.Context, ID uuid.UUID) (*Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *database) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.ID = id
	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return db.Get(ctx, id)
}

func (db *database) Delete(ctx context.Context, ID uuid.UUID) error {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}

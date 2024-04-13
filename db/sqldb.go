package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
	"github.com/mrinalwahal/service/pkg/middleware"
	"gorm.io/gorm"
)

type SQLDBConfig struct {

	// Database connection.
	// The connection should already be open.
	//
	// This field is mandatory.
	DB *gorm.DB
}

func NewSQLDB(config *SQLDBConfig) DB {
	db := sqldb{
		conn: config.DB,
	}

	return &db
}

// sqldb is the database layer implementation of an SQL/Relational type database.
//
// For example, MySQL, PostgreSQL, SQLite, etc.
//
// It implements the DB interface.
type sqldb struct {

	//	Database Connection
	conn *gorm.DB
}

// Create operation creates a new record in the database.
func (db *sqldb) Create(ctx context.Context, options *CreateOptions) (*model.Record, error) {
	txn := db.conn.WithContext(ctx)
	if options == nil {
		return nil, ErrInvalidOptions
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	//
	// This method has no Row Level Security (RLS) checks.
	//

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

// List operation fetches a list of records from the database.
func (db *sqldb) List(ctx context.Context, options *ListOptions) ([]*model.Record, error) {
	txn := db.conn.WithContext(ctx)
	if options == nil {
		options = &ListOptions{}
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can list it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

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

	if result := query.Find(&payload); result.Error != nil {
		return nil, result.Error
	}
	return payload, nil
}

// Get operation fetches a record from the database.
func (db *sqldb) Get(ctx context.Context, ID uuid.UUID) (*model.Record, error) {
	txn := db.conn.WithContext(ctx)
	if ID == uuid.Nil {
		return nil, ErrInvalidOptions
	}

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can get it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

	var payload model.Record
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

// Update operation updates a record in the database.
func (db *sqldb) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*model.Record, error) {
	txn := db.conn.WithContext(ctx)
	if id == uuid.Nil || options == nil {
		return nil, ErrInvalidOptions
	}
	if err := options.validate(); err != nil {
		return nil, err
	}

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can update it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

	var payload model.Record
	payload.ID = id
	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return db.Get(ctx, id)
}

// Delete operation deletes a record from the database.
func (db *sqldb) Delete(ctx context.Context, ID uuid.UUID) error {
	txn := db.conn.WithContext(ctx)
	if ID == uuid.Nil {
		return ErrInvalidOptions
	}

	// If the request context contains JWT claims, apply Row Level Security (RLS) checks.
	claims, exists := ctx.Value(middleware.XJWTClaims).(middleware.JWTClaims)
	if exists {

		// 1. Only the user who created the record can delete it.
		txn = txn.Where(&model.Record{
			UserID: claims.XUserID,
		})
	}

	var payload model.Record
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}

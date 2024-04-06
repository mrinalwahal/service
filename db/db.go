package db

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	slogGorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"
)

type DB interface {
	Create(context.Context, *CreateOptions) (*Record, error)
	Get(context.Context, uuid.UUID) (*Record, error)
	List(context.Context, *ListOptions) ([]*Record, error)
	Update(context.Context, uuid.UUID, *UpdateOptions) (*Record, error)
	Delete(context.Context, uuid.UUID) error
}

type Config struct {

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

func NewDB(config *Config) (DB, error) {

	logger := config.Logger
	if logger == nil {
		logger = slog.Default()
	}

	//	Setup the gorm logger.
	handler := logger.With("layer", "database").Handler()
	gormLogger := slogGorm.New(
		slogGorm.WithHandler(handler), // since v1.3.0
		slogGorm.WithTraceAll(),       // trace all messages
	)

	conn, err := gorm.Open(config.Dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	db := Database{
		conn: conn,
	}

	return &db, nil
}

type Database struct {

	//	Database Connection
	conn *gorm.DB
}

func (db *Database) Create(ctx context.Context, options *CreateOptions) (*Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.Title = options.Title

	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *Database) Get(ctx context.Context, ID uuid.UUID) (*Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func (db *Database) List(ctx context.Context, options *ListOptions) ([]*Record, error) {
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

func (db *Database) Update(ctx context.Context, id uuid.UUID, options *UpdateOptions) (*Record, error) {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.ID = id
	if result := txn.Model(&payload).Updates(options); result.Error != nil {
		return nil, result.Error
	}
	return db.Get(ctx, id)
}

func (db *Database) Delete(ctx context.Context, ID uuid.UUID) error {
	txn := db.conn.WithContext(ctx)

	var payload Record
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}

func create(txn *gorm.DB, options *CreateOptions) (*Record, error) {
	var payload Record
	payload.Title = options.Title

	result := txn.Create(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func get(txn *gorm.DB, ID uuid.UUID) (*Record, error) {
	var payload Record
	payload.ID = ID
	result := txn.First(&payload)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payload, nil
}

func list(txn *gorm.DB, options *ListOptions) ([]*Record, error) {
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

func update(txn *gorm.DB, id uuid.UUID, options *UpdateOptions) error {
	var payload Record
	payload.ID = id
	result := txn.Model(&payload).Updates(options)
	return result.Error
}

func delete(txn *gorm.DB, ID uuid.UUID) error {
	var payload Record
	payload.ID = ID
	result := txn.Delete(&payload)
	return result.Error
}

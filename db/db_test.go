package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Temporary environment that contains all the configuration required by our tests.
type environment struct {

	// Test database connection.
	conn *gorm.DB
}

// Setup the test environment.
func initialize(t *testing.T) *environment {

	// Open an in-memory database connection with SQLite.
	conn, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open the database connection: %v", err)
	}

	// Migrate the schema.
	if err := conn.AutoMigrate(&model.Record{}); err != nil {
		t.Fatalf("failed to migrate the schema: %v", err)
	}

	// Cleanup the environment after the test is complete.
	t.Cleanup(func() {

		// Close the connection.
		sqlDB, err := conn.DB()
		if err != nil {
			t.Fatalf("failed to get the database connection: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("failed to close the database connection: %v", err)
		}
	})

	return &environment{
		conn: conn,
	}
}

func Test_Database_Create(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	t.Run("create record", func(t *testing.T) {

		options := &CreateOptions{
			Title:  "Test Record",
			UserID: uuid.New(),
		}

		record, err := db.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		if record.ID == uuid.Nil {
			t.Fatalf("expected record ID to be generated automatically, got empty UUID")
		}

		if record.Title != options.Title {
			t.Fatalf("expected record title to be 'Test Record', got '%s'", record.Title)
		}
	})

	t.Run("empty title", func(t *testing.T) {

		options := &CreateOptions{
			Title:  "",
			UserID: uuid.New(),
		}

		_, err := db.Create(context.Background(), options)
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
	})

	t.Run("generate UUID of a new record automatically", func(t *testing.T) {

		options := &CreateOptions{
			Title:  "Test Record",
			UserID: uuid.New(),
		}

		record, err := db.Create(context.Background(), options)
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Check if the response contains a valid UUID and correct title.
		if record.ID == uuid.Nil {
			t.Fatalf("expected record ID to be generated automatically, got empty UUID")
		}
	})
}

func Test_Database_List(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with some records.
	for i := 0; i < 2; i++ {
		_, err := db.Create(context.Background(), &CreateOptions{
			Title:  fmt.Sprintf("model.Record %d", i),
			UserID: uuid.New(),
		})
		if err != nil {
			t.Fatalf("failed to seed the database: %v", err)
		}
	}

	type args struct {
		ctx     context.Context
		options *ListOptions
	}
	tests := []struct {

		// The name of our test.
		// This will be used to identify the test in the output.
		//
		// Example: "list all records"
		name string

		// The arguments that we will pass to the function.
		//
		// Example: context.Background(), &CreateOptions{Title: "Test model.Record"}
		args args

		// The validation function that will be used to validate the output.
		validation func([]*model.Record) error

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "list records",
			args: args{
				ctx:     context.Background(),
				options: &ListOptions{},
			},
			validation: func(records []*model.Record) error {
				if len(records) < 1 {
					return fmt.Errorf("expected at least 1 seed record, got %d", len(records))
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.List(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("database.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validation != nil && tt.validation(got) != nil {
				t.Errorf("database.List() = %v, validation produced = %v", got, tt.validation(got))
			}
		})
	}
}

func Test_Database_Get(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	seed, err := db.Create(context.Background(), &CreateOptions{
		Title:  "Test model.Record",
		UserID: uuid.New(),
	})
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	tests := []struct {

		// The name of our test.
		// This will be used to identify the test in the output.
		//
		// Example: "list all records"
		name string

		// The arguments that we will pass to the function.
		//
		// Example: context.Background(), &CreateOptions{Title: "Test model.Record"}
		args args

		// The validation function that will be used to validate the output.
		validation func(*model.Record) error

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "Get seed record",
			args: args{
				ctx: context.Background(),
				ID:  seed.ID,
			},
			validation: func(r *model.Record) error {
				if r.ID != seed.ID {
					return fmt.Errorf("expected retrieved record to equal seed, got = %v", r)
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Get(tt.args.ctx, tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("database.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validation != nil && tt.validation(got) != nil {
				t.Errorf("database.Get() = %v, validation produced = %v", got, tt.validation(got))
			}
		})
	}
}

func Test_Database_Update(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	seed, err := db.Create(context.Background(), &CreateOptions{
		Title:  "Test model.Record",
		UserID: uuid.New(),
	})
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	type args struct {
		ctx     context.Context
		id      uuid.UUID
		options *UpdateOptions
	}
	tests := []struct {

		// The name of our test.
		// This will be used to identify the test in the output.
		//
		// Example: "list all records"
		name string

		// The arguments that we will pass to the function.
		//
		// Example: context.Background(), &CreateOptions{Title: "Test model.Record"}
		args args

		// The validation function that will be used to validate the output.
		validation func(*model.Record) error

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "Update seed record",
			args: args{
				ctx: context.Background(),
				id:  seed.ID,
				options: &UpdateOptions{
					Title: "Updated Title",
				},
			},
			validation: func(r *model.Record) error {
				if r.Title != "Updated Title" {
					return fmt.Errorf("expected updated record title to be 'Updated Title', got '%s'", r.Title)
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Update(tt.args.ctx, tt.args.id, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("database.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validation != nil && tt.validation(got) != nil {
				t.Errorf("database.Update() = %v, validation produced = %v", got, tt.validation(got))
			}
		})
	}
}

func Test_Database_Delete(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	seed, err := db.Create(context.Background(), &CreateOptions{
		Title:  "Test model.Record",
		UserID: uuid.New(),
	})
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	tests := []struct {

		// The name of our test.
		// This will be used to identify the test in the output.
		//
		// Example: "list all records"
		name string

		// The arguments that we will pass to the function.
		//
		// Example: context.Background(), &CreateOptions{Title: "Test model.Record"}
		args args

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "Delete seed record",
			args: args{
				ctx: context.Background(),
				ID:  seed.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Delete(tt.args.ctx, tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("database.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Temporary environment that contains all the configuration required by our tests.
type environment struct {
	conn *gorm.DB
}

// Setup the test environment.
func setup(t *testing.T) *environment {

	// Open a test database connection with SQLite.
	conn, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open the database connection: %v", err)
	}

	// Migrate the schema.
	if err := conn.AutoMigrate(&Record{}); err != nil {
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

func Test_database_Create(t *testing.T) {

	// Setup the test environment.
	environment := setup(t)

	type fields struct {
		conn *gorm.DB
	}
	type args struct {
		ctx     context.Context
		options *CreateOptions
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		validation func(*Record) error
		wantErr    bool
	}{
		{
			name: "Create a record",
			fields: fields{
				conn: environment.conn,
			},
			args: args{
				ctx: context.Background(),
				options: &CreateOptions{
					Title: "Test Title",
				},
			},
			validation: func(r *Record) error {
				if r.Title != "Test Title" {
					return fmt.Errorf("expected record title to be 'Test Title', got '%s'", r.Title)
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "Empty title",
			fields: fields{
				conn: environment.conn,
			},
			args: args{
				ctx: context.Background(),
				options: &CreateOptions{
					Title: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Generate UUID of a new record automatically",
			fields: fields{
				conn: environment.conn,
			},
			args: args{
				ctx: context.Background(),
				options: &CreateOptions{
					Title: "Test Title",
				},
			},
			validation: func(r *Record) error {
				if len(r.ID.String()) == 0 {
					return fmt.Errorf("expected record ID to be generated automatically, got empty UUID")
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &database{
				conn: tt.fields.conn,
			}
			got, err := db.Create(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("database.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validation != nil && tt.validation(got) != nil {
				t.Errorf("database.Create() = %v, validation produced = %v", got, tt.validation(got))
			}
		})
	}
}

func Test_database_List(t *testing.T) {

	// Setup the test environment.
	environment := setup(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with some records.
	for i := 0; i < 2; i++ {
		_, err := db.Create(context.Background(), &CreateOptions{
			Title: fmt.Sprintf("Record %d", i),
		})
		if err != nil {
			t.Fatalf("failed to seed the database: %v", err)
		}
	}

	type fields struct {
		conn *gorm.DB
	}
	type args struct {
		ctx     context.Context
		options *ListOptions
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		validation func([]*Record) error
		wantErr    bool
	}{
		{
			name: "List records",
			fields: fields{
				conn: environment.conn,
			},
			args: args{
				ctx:     context.Background(),
				options: &ListOptions{},
			},
			validation: func(records []*Record) error {
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

func Test_database_Get(t *testing.T) {

	// Setup the test environment.
	environment := setup(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	seed, err := db.Create(context.Background(), &CreateOptions{
		Title: "Test Record",
	})
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	type fields struct {
		conn *gorm.DB
	}
	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		validation func(*Record) error
		wantErr    bool
	}{
		{
			name: "Get seed record",
			args: args{
				ctx: context.Background(),
				ID:  seed.ID,
			},
			validation: func(r *Record) error {
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

func Test_database_Update(t *testing.T) {

	// Setup the test environment.
	environment := setup(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	seed, err := db.Create(context.Background(), &CreateOptions{
		Title: "Test Record",
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
		name       string
		args       args
		validation func(*Record) error
		wantErr    bool
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
			validation: func(r *Record) error {
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

func Test_database_Delete(t *testing.T) {

	// Setup the test environment.
	environment := setup(t)

	// Initialize the database.
	db := &database{
		conn: environment.conn,
	}

	// Seed the database with sample records.
	seed, err := db.Create(context.Background(), &CreateOptions{
		Title: "Test Record",
	})
	if err != nil {
		t.Fatalf("failed to seed the database: %v", err)
	}

	type args struct {
		ctx context.Context
		ID  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
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

package service

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/model"
	"go.uber.org/mock/gomock"
)

// Contains all the configuration required by our tests.
type testconfig struct {

	// Mock database layer.
	db *db.MockDB

	// Test log.
	log *slog.Logger
}

// Setup the test environment.
func configure(t *testing.T) *testconfig {

	// Get the mock database layer.
	db := db.NewMockDB(gomock.NewController(t))
	return &testconfig{
		db:  db,
		log: slog.Default(),
	}
}

func Test_Service_Create(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	t.Run("create record with nil options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Create(context.Background(), nil)
		if err == nil || err != ErrOptionsNotFound {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create record with invalid options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Create(context.Background(), &CreateOptions{
			Title: "",
		})
		if err == nil {
			t.Errorf("service.Create() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("create record with valid options", func(t *testing.T) {

		record := model.Record{
			Title: "Test Record",
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Record{
			Base: model.Base{
				ID: uuid.New(),
			},
			Title: record.Title,
		}, nil).Times(1)

		got, err := s.Create(context.Background(), &CreateOptions{
			Title: record.Title,
		})
		if err != nil {
			t.Errorf("service.Create() error = %v, wantErr %v", err, false)
		}
		if got.ID == uuid.Nil {
			t.Errorf("service.Create() = %v, want a valid UUID", got.ID)
		}
		if got.Title != record.Title {
			t.Errorf("service.Create() = %v, want %v", got.Title, record.Title)
		}
	})
}

func Test_Service_List(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
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

		// The expectation that we will set on the mock database layer.
		expectation *gomock.Call

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
			expectation: config.db.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.Record{
				{
					Title: "Record 1",
				},
				{
					Title: "Record 2",
				},
			}, nil),
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

			// Set the expectation.
			tt.expectation.Times(1)

			got, err := s.List(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validation != nil && tt.validation(got) != nil {
				t.Errorf("service.List() = %v, validation produced = %v", got, tt.validation(got))
			}
		})
	}
}

func Test_Service_Get(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	// Sample record UUID.
	id := uuid.New()

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

		// The expectation that we will set on the mock database layer.
		expectation *gomock.Call

		// The validation function that will be used to validate the output.
		validation func(*model.Record) error

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "get seed record",
			args: args{
				ctx: context.Background(),
				ID:  id,
			},
			expectation: config.db.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&model.Record{
				Base: model.Base{
					ID: id,
				},
				Title: "Test Record",
			}, nil),
			validation: func(record *model.Record) error {
				if record.ID != id {
					return fmt.Errorf("expected retrieved record to equal seed, got = %v", record.ID)
				}
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Set the expectation.
			tt.expectation.Times(1)

			got, err := s.Get(tt.args.ctx, tt.args.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validation != nil && tt.validation(got) != nil {
				t.Errorf("service.List() = %v, validation produced = %v", got, tt.validation(got))
			}
		})
	}
}

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
	"gorm.io/gorm"
)

// Temporary environment that contains all the configuration required by our tests.
type environment struct {

	// Mock database layer.
	db *db.MockDB

	// Test logger.
	logger *slog.Logger
}

// Setup the test environment.
func initialize(t *testing.T) *environment {

	// Get the mock database layer.
	db := db.NewMockDB(gomock.NewController(t))
	return &environment{
		db:     db,
		logger: slog.Default(),
	}
}

func Test_Service_Create(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the service.
	s := &service{
		db:     environment.db,
		logger: environment.logger,
	}

	type args struct {
		ctx     context.Context
		options *CreateOptions
	}
	tests := []struct {

		// The name of our test.
		// This will be used to identify the test in the output.
		//
		// Example: "create record"
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
			name: "create record",
			args: args{
				ctx: context.Background(),
				options: &CreateOptions{
					Title: "Test model.Record",
				},
			},
			expectation: environment.db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Record{
				Title: "Test model.Record",
			}, nil),
			validation: func(r *model.Record) error {
				if r.Title != "Test model.Record" {
					return fmt.Errorf("expected record title to be 'Test model.Record', got '%s'", r.Title)
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "empty title",
			args: args{
				ctx: context.Background(),
				options: &CreateOptions{
					Title: "",
				},
			},
			expectation: environment.db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrInvalidValue),
			wantErr:     true,
		},
		{
			name: "generate UUID of a new record automatically",
			args: args{
				ctx: context.Background(),
				options: &CreateOptions{
					Title: "Test model.Record",
				},
			},
			expectation: environment.db.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Record{
				Base: model.Base{
					ID: uuid.New(),
				},
				Title: "Test model.Record",
			}, nil),
			validation: func(r *model.Record) error {
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

			// Set the expectation.
			tt.expectation.Times(1)

			got, err := s.Create(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.validation != nil && tt.validation(got) != nil {
				t.Errorf("service.Create() = %v, validation produced = %v", got, tt.validation(got))
			}
		})
	}
}

func Test_Service_List(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the service.
	s := &service{
		db:     environment.db,
		logger: environment.logger,
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
			expectation: environment.db.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.Record{
				{
					Title: "model.Record 1",
				},
				{
					Title: "model.Record 2",
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

	// Setup the test environment.
	environment := initialize(t)

	// Initialize the service.
	s := &service{
		db:     environment.db,
		logger: environment.logger,
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
			expectation: environment.db.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&model.Record{
				Base: model.Base{
					ID: id,
				},
				Title: "Test model.Record",
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

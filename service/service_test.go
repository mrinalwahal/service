package service

import (
	"context"
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

func Test_NewService(t *testing.T) {

	t.Run("nil config", func(t *testing.T) {

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewService() did not panic")
			}
		}()

		// Initialize the service.
		NewService(nil)
	})

	t.Run("valid config w/ db", func(t *testing.T) {

		// Get the mock database layer.
		db := db.NewMockDB(gomock.NewController(t))

		// Initialize the service.
		s := NewService(&Config{
			DB: db,
		})

		if s == nil {
			t.Errorf("NewService() = %v, want a valid service", s)
		}
	})

	t.Run("valid config w/ db and logger", func(t *testing.T) {

		// Get the mock database layer.
		db := db.NewMockDB(gomock.NewController(t))

		// Initialize the service.
		s := NewService(&Config{
			DB:     db,
			Logger: slog.Default(),
		})

		if s == nil {
			t.Errorf("NewService() = %v, want a valid service", s)
		}
	})
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
		if err == nil || err != ErrInvalidOptions {
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
			Title:  record.Title,
			UserID: uuid.New(),
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

	t.Run("list records with nil options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().List(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.List(context.Background(), nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.List() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("list records with invalid options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().List(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.List(context.Background(), &ListOptions{
			Skip:  -1,
			Limit: -1,
		})
		if err == nil {
			t.Errorf("service.List() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("list records with valid options", func(t *testing.T) {

		records := []*model.Record{
			{
				Base: model.Base{
					ID: uuid.New(),
				},
				Title: "Test Record",
			},
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().List(gomock.Any(), gomock.Any()).Return(records, nil).Times(1)

		got, err := s.List(context.Background(), &ListOptions{
			Skip:  0,
			Limit: 10,
		})
		if err != nil {
			t.Errorf("service.List() error = %v, wantErr %v", err, false)
		}
		if len(got) != len(records) {
			t.Errorf("service.List() = %v, want %v", len(got), len(records))
		}
	})
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

	t.Run("get record with invalid ID", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Get(context.Background(), uuid.Nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Get() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("get record with valid ID", func(t *testing.T) {

		record := model.Record{
			Base: model.Base{
				ID: id,
			},
			Title: "Test Record",
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().Get(gomock.Any(), id).Return(&record, nil).Times(1)

		got, err := s.Get(context.Background(), id)
		if err != nil {
			t.Errorf("service.Get() error = %v, wantErr %v", err, false)
		}
		if got.ID != id {
			t.Errorf("service.Get() = %v, want %v", got.ID, id)
		}
		if got.Title != record.Title {
			t.Errorf("service.Get() = %v, want %v", got.Title, record.Title)
		}
	})
}

func Test_Service_Update(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	// Sample record UUID.
	id := uuid.New()

	t.Run("update record with invalid ID", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Update(context.Background(), uuid.Nil, &UpdateOptions{
			Title: "Test Record",
		})
		if err == nil || err != ErrInvalidRecordID {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update record with nil options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Update(context.Background(), id, nil)
		if err == nil || err != ErrInvalidOptions {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update record with invalid options", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		_, err := s.Update(context.Background(), id, &UpdateOptions{
			Title: "",
		})
		if err == nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("update record with valid options", func(t *testing.T) {

		record := model.Record{
			Base: model.Base{
				ID: id,
			},
			Title: "Test Record",
		}

		// Set the expectation at the database layer.
		config.db.EXPECT().Update(gomock.Any(), id, gomock.Any()).Return(&record, nil).Times(1)

		got, err := s.Update(context.Background(), id, &UpdateOptions{
			Title: "Updated Record",
		})
		if err != nil {
			t.Errorf("service.Update() error = %v, wantErr %v", err, false)
		}
		if got.ID != id {
			t.Errorf("service.Update() = %v, want %v", got.ID, id)
		}
		if got.Title != record.Title {
			t.Errorf("service.Update() = %v, want %v", got.Title, record.Title)
		}
	})
}

func Test_Service_Delete(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	// Initialize the service.
	s := &service{
		db:     config.db,
		logger: config.log,
	}

	// Sample record UUID.
	id := uuid.New()

	t.Run("delete record with invalid ID", func(t *testing.T) {

		// Make sure the database layer is not expecting a call.
		config.db.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)

		err := s.Delete(context.Background(), uuid.Nil)
		if err == nil || err != ErrInvalidRecordID {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, true)
		}
	})

	t.Run("delete record with valid ID", func(t *testing.T) {

		// Set the expectation at the database layer.
		config.db.EXPECT().Delete(gomock.Any(), id).Return(nil).Times(1)

		err := s.Delete(context.Background(), id)
		if err != nil {
			t.Errorf("service.Delete() error = %v, wantErr %v", err, false)
		}
	})
}

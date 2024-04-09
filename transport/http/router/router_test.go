package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrinalwahal/service/db"
	"github.com/mrinalwahal/service/model"
	"github.com/mrinalwahal/service/service"
	"github.com/mrinalwahal/service/transport/http/handler"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testconfig contains all the configuration that is required by our tests.
type testconfig struct {

	// Logger instance.
	log *slog.Logger

	// Service layer.
	service service.Service
}

// configure configures a suitable and reliable environment for the tests.
func configure(t *testing.T) *testconfig {

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

	// Initialize the database layer.
	db := db.NewDB(&db.Config{
		DB: conn,
	})

	// Initialize the service.
	service := service.NewService(&service.Config{
		DB:     db,
		Logger: slog.Default(),
	})

	return &testconfig{
		service: service,
		log:     slog.Default(),
	}
}

func Test_Router(t *testing.T) {

	// Configure the test environment.
	config := configure(t)

	t.Run("request to create record w/ no body", func(t *testing.T) {

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodPost, "/v1", nil)
		response := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(response, request)

		// Check the response status code.
		if response.Code != http.StatusBadRequest {
			t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, response.Code)
		}
	})

	t.Run("request to create record w/ invalid body", func(t *testing.T) {

		// Prepare a body with invalid JSON.
		body, err := json.Marshal(map[string]interface{}{
			"invalid": "json",
		})
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodPost, "/v1", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusBadRequest {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("request to create record w/ valid body", func(t *testing.T) {

		// Prepare a body with invalid JSON.
		body, err := json.Marshal(handler.CreateOptions{
			Title: "test",
		})
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodPost, "/v1", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusCreated {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusCreated, recorder.Code)
		}
	})

	t.Run("request to get record w/ invalid id", func(t *testing.T) {

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodGet, "/v1/invalid-id", nil)
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusBadRequest {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("request to get record w/ valid id", func(t *testing.T) {

		// Create a record.
		record, err := config.service.Create(context.Background(), &service.CreateOptions{
			Title: "test",
		})
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/%s", record.ID), nil)
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusOK {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("request to list records", func(t *testing.T) {

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodGet, "/v1", nil)
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusOK {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		// Check that the returned data in the response is a JSON array.
		var response handler.Response
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal the response body: %v", err)
		}

		if response.Data == nil {
			t.Fatalf("expected response data to be a JSON array, got nil")
		}
	})

	t.Run("request to update record w/ invalid id", func(t *testing.T) {

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodPatch, "/v1/invalid-id", nil)
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusBadRequest {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("request to update record w/ valid id", func(t *testing.T) {

		// Create a record.
		record, err := config.service.Create(context.Background(), &service.CreateOptions{
			Title: "test",
		})
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Prepare the body.
		body, err := json.Marshal(handler.UpdateOptions{
			Title: "updated",
		})
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/%s", record.ID), bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusOK {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		// Validate the title of the updated record.
		var response handler.Response
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal the response body: %v", err)
		}

		if response.Data == nil {
			t.Fatalf("expected response data to be a JSON object, got nil")
		}

		data, ok := response.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("expected response data to be a JSON object, got %T", response.Data)
		}

		if data["title"] != "updated" {
			t.Fatalf("expected title to be 'updated', got %s", data["title"])
		}
	})

	t.Run("request to delete record w/ invalid id", func(t *testing.T) {

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodDelete, "/v1/invalid-id", nil)
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusBadRequest {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("request to delete record w/ valid id", func(t *testing.T) {

		// Create a record.
		record, err := config.service.Create(context.Background(), &service.CreateOptions{
			Title: "test",
		})
		if err != nil {
			t.Fatalf("failed to create a record: %v", err)
		}

		// Prepare the request and response recorder.
		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/%s", record.ID), nil)
		recorder := httptest.NewRecorder()

		// Prepare the router.
		router := NewHTTPRouter(&HTTPRouterConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Serve the request.
		router.ServeHTTP(recorder, request)

		// Check the response status code.
		if recorder.Code != http.StatusOK {
			t.Logf("got response body = %v", recorder.Body.String())
			t.Fatalf("expected status code %d, got %d", http.StatusOK, recorder.Code)
		}

		// Try to fetch the deleted record and ensure it doesn't exist.
		_, err = config.service.Get(context.Background(), record.ID)
		if err == nil {
			t.Fatal("expected to get an error, got nil")
		}
	})
}

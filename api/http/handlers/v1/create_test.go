package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/pkg/middleware"
	"github.com/mrinalwahal/service/service"
	"go.uber.org/mock/gomock"
)

// Temporary environment that contains all the configuration required by our tests.
type environment struct {

	// Mock service layer.
	service *service.MockService

	// Test logger.
	logger *slog.Logger
}

// Setup the test environment.
func initialize(t *testing.T) *environment {

	// Get the mock service layer.
	service := service.NewMockService(gomock.NewController(t))
	return &environment{
		service: service,
		logger:  slog.Default(),
	}
}

func TestCreateHandler_ServeHTTP(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	t.Run("create record without sending UserID", func(t *testing.T) {

		// Initialize the handler.
		h := &CreateHandler{
			service: environment.service,
			log:     environment.logger,
		}

		// Prepare the request body.
		body, err := json.Marshal(CreateOptions{
			Title: "Test",
		})
		if err != nil {
			t.Fatal(err)
		}

		// Create a new request.
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		// Create a new response recorder.
		w := httptest.NewRecorder()

		// Set the expectation for the service layer.
		environment.service.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		// Serve the request.
		h.ServeHTTP(w, req)

		// Validate the status code.
		if status := w.Code; status != http.StatusUnauthorized {
			t.Logf("Response: %s", w.Body.String())
			t.Errorf("CreateHandler.ServeHTTP() = %v, want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("create record with nil options", func(t *testing.T) {

		// Initialize the handler.
		h := &CreateHandler{
			service: environment.service,
			log:     environment.logger,
		}

		// Create a new request.
		req := httptest.NewRequest(http.MethodPost, "/", nil)

		// Set random UserID in the request context.
		req = req.WithContext(context.WithValue(req.Context(), middleware.XJWTClaims, JWTClaims{
			UserID: uuid.New(),
		}))

		// Create a new response recorder.
		w := httptest.NewRecorder()

		// Set the expectation for the service layer.
		environment.service.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		// Serve the request.
		h.ServeHTTP(w, req)

		// Validate the status code.
		if status := w.Code; status != http.StatusBadRequest {
			t.Logf("Response: %s", w.Body.String())
			t.Errorf("CreateHandler.ServeHTTP() = %v, want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("create record with empty title", func(t *testing.T) {

		// Initialize the handler.
		h := &CreateHandler{
			service: environment.service,
			log:     environment.logger,
		}

		// Prepare the request body.
		body, err := json.Marshal(CreateOptions{
			Title: "",
		})
		if err != nil {
			t.Fatal(err)
		}

		// Create a new request.
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		// Set random UserID in the request context.
		req = req.WithContext(context.WithValue(req.Context(), middleware.XJWTClaims, JWTClaims{
			UserID: uuid.New(),
		}))

		// Create a new response recorder.
		w := httptest.NewRecorder()

		// Set the expectation for the service layer.
		environment.service.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		// Serve the request.
		h.ServeHTTP(w, req)

		// Validate the status code.
		if status := w.Code; status != http.StatusBadRequest {
			t.Logf("Response: %s", w.Body.String())
			t.Errorf("CreateHandler.ServeHTTP() = %v, want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("create record with valid options", func(t *testing.T) {

		// Initialize the handler.
		h := &CreateHandler{
			service: environment.service,
			log:     environment.logger,
		}

		// Prepare the request body.
		body, err := json.Marshal(CreateOptions{
			Title: "Test",
		})
		if err != nil {
			t.Fatal(err)
		}

		// Create a new request.
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

		// Set random UserID in the request context.
		req = req.WithContext(context.WithValue(req.Context(), middleware.XJWTClaims, JWTClaims{
			UserID: uuid.New(),
		}))

		// Create a new response recorder.
		w := httptest.NewRecorder()

		// Set the expectation for the service layer.
		environment.service.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		// Serve the request.
		h.ServeHTTP(w, req)

		// Validate the status code.
		if status := w.Code; status != http.StatusCreated {
			t.Logf("Response: %s", w.Body.String())
			t.Errorf("CreateHandler.ServeHTTP() = %v, want %v", status, http.StatusCreated)
		}
	})
}

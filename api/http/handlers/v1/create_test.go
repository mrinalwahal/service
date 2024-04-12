package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
	"github.com/mrinalwahal/service/pkg/middleware"
	"github.com/mrinalwahal/service/service"
	"go.uber.org/mock/gomock"
)

// Contains all the configuration required by our tests.
type testconfig struct {

	// Mock service layer.
	service *service.MockService

	// Test log.
	log *slog.Logger
}

// Setup the test environment.
func configure(t *testing.T) *testconfig {

	// Get the mock service layer.
	service := service.NewMockService(gomock.NewController(t))
	return &testconfig{
		service: service,
		log:     slog.Default(),
	}
}

func TestCreateHandler_ServeHTTP(t *testing.T) {

	// Setup the test config.
	config := configure(t)

	t.Run("create w/ invalid options", func(t *testing.T) {

		// Create the handler.
		handler := NewCreateHandler(&CreateHandlerConfig{
			Service: config.service,
			Logger:  config.log,
		})

		// Initialize test request and response recorder.
		r := httptest.NewRequest(http.MethodPost, "/v1/records", nil)
		w := httptest.NewRecorder()

		// The service layer should ideally not be expecting any calls to reach it.
		config.service.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)

		// Serve the request.
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("create w/ valid options but w/o jwt claims", func(t *testing.T) {

		// Create the handler.
		handler := NewCreateHandler(&CreateHandlerConfig{
			Service: config.service,
			Logger:  config.log,
		})

		body, err := json.Marshal(CreateOptions{
			Title: "Test Record",
		})
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Initialize test request and response recorder.
		r := httptest.NewRequest(http.MethodPost, "/v1/records", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// The service layer should ideally return an error because the JWT claims are missing.
		config.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, service.ErrInvalidOptions)

		// Serve the request.
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("create w/ valid options and jwt claims", func(t *testing.T) {

		// Create the handler.
		handler := NewCreateHandler(&CreateHandlerConfig{
			Service: config.service,
			Logger:  config.log,
		})

		options := CreateOptions{
			Title: "Test Record",
		}
		body, err := json.Marshal(options)
		if err != nil {
			t.Fatalf("failed to marshal the dummy body for request: %v", err)
		}

		// Initialize test request and response recorder.
		r := httptest.NewRequest(http.MethodPost, "/v1/records", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		// Set the JWT claims in the request context.
		id := uuid.New()
		r = r.WithContext(context.WithValue(r.Context(), middleware.XJWTClaims, jwt.MapClaims{
			"x-user-id": id,
		}))

		// The service layer should ideally return a record.
		config.service.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.Record{
			Base: model.Base{
				ID: uuid.New(),
			},
			Title:  options.Title,
			UserID: id,
		}, nil)

		// Serve the request.
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})
}

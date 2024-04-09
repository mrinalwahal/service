package v1

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

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

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {

		// The name of our test.
		// This will be used to identify the test in the output.
		//
		// Example: "create record"
		name string

		// The arguments that we will pass to the function.
		//
		// Example: `w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/", nil)`
		args args

		// The expectation that we will set on the mock database layer.
		expectation *gomock.Call

		// The status code we expect in response.
		//
		// Example: http.StatusOK
		want int
	}{
		{
			name: "create record",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"title":"test"}`)),
			},
			expectation: environment.service.EXPECT().Create(gomock.Any(), gomock.Any()),
			want:        http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateHandler{
				service: environment.service,
				log:     environment.logger,
			}

			// Set the expectation.
			tt.expectation.Times(1)

			h.ServeHTTP(tt.args.w, tt.args.r)

			// Validate the status code.
			if status := tt.args.w.(*httptest.ResponseRecorder).Code; status != tt.want {
				t.Log(tt.args.w.(*httptest.ResponseRecorder).Body.String())
				t.Errorf("CreateHandler.ServeHTTP() = %v, want %v", status, tt.want)
			}
		})
	}
}

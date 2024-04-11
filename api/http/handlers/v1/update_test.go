package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
	"github.com/mrinalwahal/service/service"
	"go.uber.org/mock/gomock"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {

	// Setup the test environment.
	environment := initialize(t)

	// Test UUID of the record.
	recordID := uuid.New()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {

		// The name of our test.
		// This will be used to identify the test in the output.
		//
		// Example: "get a record"
		name string

		// The arguments that we will pass to the function.
		//
		// Example: `w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/", nil)`
		args args

		// The expectation that we will set on the mock database layer.
		expectation *gomock.Call

		// The validation function that will be used to validate the output.
		validation func(*Response) error

		// The status code we expect in response.
		//
		// Example: http.StatusOK
		wantStatus int

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "update record succesfully",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s", recordID.String()), bytes.NewBufferString(`{"title": "Updated Title"}`))
					req.SetPathValue("id", recordID.String())
					return req
				}(),
			},
			expectation: environment.service.EXPECT().Update(gomock.Any(), recordID, &service.UpdateOptions{
				Title: "Updated Title",
			}, gomock.Any()).Return(&model.Record{
				Title: "Updated Title",
			}, nil),
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "return invalid title after updating record",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s", recordID.String()), bytes.NewBufferString(`{"title": "Updated Title"}`))
					req.SetPathValue("id", recordID.String())
					return req
				}(),
			},
			expectation: environment.service.EXPECT().Update(gomock.Any(), recordID, &service.UpdateOptions{
				Title: "Updated Title",
			}, gomock.Any()).Return(&model.Record{
				Title: "Wrong Title",
			}, nil),
			validation: func(r *Response) error {
				if r.Message != "Updated title" {
					return fmt.Errorf("expected message to be 'Updated title', got %s", r.Message)
				}
				return nil
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UpdateHandler{
				service: environment.service,
				log:     environment.logger,
			}

			// Set the expectation.
			tt.expectation.Times(1)

			h.ServeHTTP(tt.args.w, tt.args.r)

			// Decode the body
			var resp Response
			if err := json.Unmarshal(tt.args.w.(*httptest.ResponseRecorder).Body.Bytes(), &resp); err != nil {
				t.Errorf("UpdateHandler.ServeHTTP() = %v", err)
			}

			// Validate the status code.
			if status := tt.args.w.(*httptest.ResponseRecorder).Code; status != tt.wantStatus {
				t.Errorf("UpdateHandler.ServeHTTP() = %v, want %v", status, tt.wantStatus)
			}

			// Run validation function.
			if tt.validation != nil {
				if err := tt.validation(&resp); (err != nil) != tt.wantErr {
					t.Errorf("UpdateHandler.ServeHTTP() = %v", err)
				}
			}
		})
	}
}

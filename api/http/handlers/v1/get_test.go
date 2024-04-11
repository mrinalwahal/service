package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/model"
	"go.uber.org/mock/gomock"
)

func TestGetHandler_ServeHTTP(t *testing.T) {

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
		want int

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "get record",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/", nil)
					req.SetPathValue("id", recordID.String())
					return req
				}(),
			},
			expectation: environment.service.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.Record{
				Base: model.Base{
					ID: recordID,
				},
				Title: "model.Record 1",
			}, nil),
			validation: func(res *Response) error {
				if res.Data == nil {
					t.Log("Response:", res)
					return fmt.Errorf("expected data to be non-nil")
				}
				return nil
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &GetHandler{
				service: environment.service,
				log:     environment.logger,
			}

			// Set the expectation.
			tt.expectation.Times(1)

			h.ServeHTTP(tt.args.w, tt.args.r)

			// Decode the body
			var resp Response
			if err := json.Unmarshal(tt.args.w.(*httptest.ResponseRecorder).Body.Bytes(), &resp); err != nil {
				t.Errorf("GetHandler.ServeHTTP() = %v", err)
			}

			// Validate the status code.
			if status := tt.args.w.(*httptest.ResponseRecorder).Code; status != tt.want {
				t.Errorf("GetHandler.ServeHTTP() = %v, want %v", status, tt.want)
			}

			// Run validation function.
			if tt.validation != nil {
				if err := tt.validation(&resp); (err != nil) != tt.wantErr {
					t.Errorf("GetHandler.ServeHTTP() = %v", err)
				}
			}
		})
	}
}

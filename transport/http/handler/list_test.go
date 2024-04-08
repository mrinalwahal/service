package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrinalwahal/service/db"
	"go.uber.org/mock/gomock"
)

func TestListHandler_ServeHTTP(t *testing.T) {

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
		// Example: "list all records"
		name string

		// The arguments that we will pass to the function.
		//
		// Example: `w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/", nil)`
		args args

		// The expectation that we will set on the mock database layer.
		expectation *gomock.Call

		// The validation function that will be used to validate the output.
		validation func(*response) error

		// The status code we expect in response.
		//
		// Example: http.StatusOK
		want int

		// Whether we expect an error or not.
		wantErr bool
	}{
		{
			name: "list all record",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			expectation: environment.service.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*db.Record{
				{
					Title: "Record 1",
				},
			}, nil),
			validation: func(r *response) error {
				if r == nil {
					return fmt.Errorf("expected a response, got nil")
				}
				records := r.Data.([]interface{})
				if len(records) < 1 {
					return fmt.Errorf("expected at least 1 record, got %d", len(records))
				}
				return nil
			},
			want: http.StatusOK,
		},
		{
			name: "list only 1 record",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"limit":1}`)),
			},
			expectation: environment.service.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*db.Record{
				{
					Title: "Record 1",
				},
			}, nil),
			validation: func(r *response) error {
				if r == nil {
					return fmt.Errorf("expected a response, got nil")
				}
				records := r.Data.([]interface{})
				if len(records) != 1 {
					return fmt.Errorf("expected only 1 record, got %d", len(records))
				}
				return nil
			},
			want: http.StatusOK,
		},
		{
			name: "return all records while requesting only 1 record",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(`{"limit":1}`)),
			},
			expectation: environment.service.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*db.Record{
				{
					Title: "Record 1",
				},
				{
					Title: "Record 2",
				},
			}, nil),
			validation: func(r *response) error {
				if r == nil {
					return fmt.Errorf("expected a response, got nil")
				}
				records := r.Data.([]interface{})
				if len(records) != 1 {
					return fmt.Errorf("expected only 1 record, got %d", len(records))
				}
				return nil
			},
			want:    http.StatusOK,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ListHandler{
				service: environment.service,
				log:     environment.logger,
			}

			// Set the expectation.
			tt.expectation.Times(1)

			h.ServeHTTP(tt.args.w, tt.args.r)

			// Decode the body
			var resp response
			if err := json.Unmarshal(tt.args.w.(*httptest.ResponseRecorder).Body.Bytes(), &resp); err != nil {
				t.Errorf("ListHandler.ServeHTTP() = %v", err)
			}

			// Validate the status code.
			if status := tt.args.w.(*httptest.ResponseRecorder).Code; status != tt.want {
				t.Log("Response:", resp)
				t.Errorf("ListHandler.ServeHTTP() = %v, want %v", status, tt.want)
			}

			// Run validation function.
			if tt.validation != nil {
				if err := tt.validation(&resp); (err != nil) != tt.wantErr {
					t.Errorf("ListHandler.ServeHTTP() = %v", err)
				}
			}
		})
	}
}

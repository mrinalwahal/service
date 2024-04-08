package handler

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/mrinalwahal/service/db"
	"go.uber.org/mock/gomock"
)

type test_environment struct {

	// Mock database layer.
	db db.DB

	// Logger instance.
	log *slog.Logger
}

func initialize(t *testing.T) *test_environment {

	// Get the mock database layer.
	db := db.NewMockDB(gomock.NewController(t))
	return &test_environment{
		db:  db,
		log: slog.Default(),
	}
}

func TestCreateHandler_ServeHTTP(t *testing.T) {

	// Initialize the environment.
	type fields struct {
		db      db.DB
		options *CreateOptions
		log     *slog.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CreateHandler{
				db:      tt.fields.db,
				options: tt.fields.options,
				log:     tt.fields.log,
			}
			h.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}

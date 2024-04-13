package db

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateOptions_validate(t *testing.T) {
	type fields struct {
		Title  string
		UserID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty title",
			fields: fields{
				Title: "",
			},
			wantErr: true,
		},
		{
			name: "invalid user id",
			fields: fields{
				UserID: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "valid options",
			fields: fields{
				Title:  "Test Record",
				UserID: uuid.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreateOptions{
				Title:  tt.fields.Title,
				UserID: tt.fields.UserID,
			}
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateOptions.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

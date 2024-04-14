package model

import "github.com/google/uuid"

type Record struct {
	Base

	// Title of the record.
	//
	// Example: "Test Record"
	//
	// It is a required field.
	Title string `json:"title" gorm:"not null;check:(length(title)>0)"`

	//	ID of the user who created the record.
	//
	//	Example: "550e8400-e29b-41d4-a716-446655440000"
	//
	//	It is a required field.
	UserID uuid.UUID `json:"user_id" gorm:"not null;type:uuid"`
}

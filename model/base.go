package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {

	// ID is the unique identifier of the object of the model.
	// It is generated automatically when the object is created. So, you cannot pass a custom value for it.
	//
	// Example: "550e8400-e29b-41d4-a716-446655440000"
	ID uuid.UUID `json:"id" gorm:"primaryKey;not null;type:uuid"`

	// CreatedAt is the time when the object was created.
	// It is set automatically when the object is created.
	//
	// Example: "2021-07-01T12:00:00Z"
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt is the time when the object was last updated.
	// It is set automatically when the object is updated.
	//
	// Example: "2021-07-01T12:00:00Z"
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// DeletedAt is the time when the object was deleted.
	// It is set automatically when the object is marked deleted.
	// Generally, used for soft deletes (marking records as deleted without actually removing them from the database).
	//
	// Example: "2021-07-01T12:00:00Z"
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// BeforeCreate hook for gorm.
// This function is called by gorm before creating a record.
//
// It performs the following operations:
//
// - Generates a new UUID for the record.
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}

package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {

	// ID is the unique identifier of the object of the model.
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`

	// CreatedAt is the time when the object was created.
	// It is set automatically when the object is created.
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt is the time when the object was last updated.
	// It is set automatically when the object is updated.
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// DeletedAt is the time when the object was deleted.
	// It is set automatically when the object is marked deleted.
	// Generally, used for soft deletes (marking records as deleted without actually removing them from the database).
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// BeforeCreate hook for gorm.
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}

type Todo struct {
	Base

	//	Title of the todo.
	Title string `json:"title" gorm:"not null"`

	//	ID of the user who created the todo.
	//UserID uuid.UUID `json:"user_id" gorm:"not null"`
}

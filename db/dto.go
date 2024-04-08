package db

import "github.com/google/uuid"

type CreateOptions struct {

	//	Title of the record.
	Title string

	// ID of the user who is creating the record.
	UserID uuid.UUID
}

// validate ascertains that CreateOptions abides by following rules:
//
// - UserID should not be nil.
// - Title should not be empty.
func (o *CreateOptions) validate() error {

	if o.UserID == uuid.Nil {
		return ErrInvalidUserID
	}
	if len(o.Title) == 0 {
		return ErrEmptyTitle
	}

	return nil
}

type ListOptions struct {

	//	Title of the record.
	Title string
	//	ID of the user who created the record.
	UserID uuid.UUID
	//	Skip for pagination.
	Skip int
	//	Limit for pagination.
	Limit int
	//	Order by field.
	OrderBy string
	//	Order by direction.
	OrderDirection string
}

type UpdateOptions struct {

	//	Title of the record.
	Title string
}

// validate ascertains that UpdateOptions abides by following rules:
//
// - Title should not be empty.
func (o *UpdateOptions) validate() error {

	if len(o.Title) == 0 {
		return ErrEmptyTitle
	}

	return nil
}

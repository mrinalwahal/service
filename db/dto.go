package db

import "github.com/google/uuid"

// Requester is the structure that holds the information of the user who sent the request.
type Requester struct {
	ID uuid.UUID
}

// CreateOptions holds the options for creating a new record.
type CreateOptions struct {

	//	Title of the record.
	Title string

	// ID of the user who is creating the record.
	UserID uuid.UUID
}

func (o *CreateOptions) validate() error {
	if o.Title == "" {
		return ErrEmptyTitle
	}
	if o.UserID == uuid.Nil {
		return ErrInvalidUserID
	}
	return nil
}

// ListOptions holds the options for listing records.
type ListOptions struct {

	//	Title of the record.
	Title string
	//	ID of the user who is created the record.
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

func (o *ListOptions) validate() error {
	if o.Skip < 0 ||
		o.Limit < 0 || o.Limit > 100 {
		return ErrInvalidFilters
	}
	return nil
}

// UpdateOptions holds the options for updating a record.
type UpdateOptions struct {

	//	Title of the record.
	Title string
}

func (o *UpdateOptions) validate() error {
	if o.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}

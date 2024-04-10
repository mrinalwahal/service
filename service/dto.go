package service

import "github.com/google/uuid"

type CreateOptions struct {

	//	Title of the record.
	Title string

	// ID of the user who is creating the record.
	UserID uuid.UUID
}

type ListOptions struct {

	//	Title of the record.
	Title string
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

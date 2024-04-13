package db

import (
	"encoding/json"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTClaims struct {
	jwt.MapClaims
	Claims  map[string]interface{} `json:"claims"`
	XUserID uuid.UUID              `json:"x-user-id"`
}

// Custom unmarshal function for JWTClaims.
func (c *JWTClaims) UnmarshalJSON(b []byte) error {
	type alias struct {
		XUserID string `json:"x-user-id"`
	}
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	userID, err := uuid.Parse(a.XUserID)
	if err != nil {
		return err
	}
	*c = JWTClaims{
		XUserID: userID,
	}
	return nil
}

// validate the JWT Claims.
func (c *JWTClaims) validate() error {
	if c.XUserID == uuid.Nil {
		return ErrInvalidUserID
	}
	return nil
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
	return nil
}

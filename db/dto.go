package db

import (
	"context"
	"encoding/json"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mrinalwahal/service/pkg/middleware"
)

type JWTClaims struct {
	XUserID uuid.UUID `json:"x-user-id"`
}

// // validate the JWT Claims.
// func (c *JWTClaims) validate() error {
// 	if c.XUserID == uuid.Nil {
// 		return ErrInvalidUserID
// 	}
// 	return nil
// }

func getClaims(ctx context.Context) (*JWTClaims, error) {
	claims, exists := ctx.Value(middleware.XJWTClaims).(jwt.MapClaims)
	if !exists {
		return nil, nil
	}
	m, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	var jwtClaims JWTClaims
	if err := json.Unmarshal(m, &jwtClaims); err != nil {
		return nil, err
	}
	return &jwtClaims, nil
}

const XUserID = "x-user-id"

const XRequestingUser = "X-Requesting-User"

// Requester is the structure that holds the information of the user who sent the request.
// type Requester struct {
// 	ID uuid.UUID `json:"id"`
// }

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
	if o.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}

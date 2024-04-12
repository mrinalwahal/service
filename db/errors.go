package db

import "fmt"

var (
	ErrInvalidOptions = fmt.Errorf("invalid options")
	ErrInvalidUserID  = fmt.Errorf("invalid user_id")
	ErrEmptyTitle     = fmt.Errorf("empty title")
	ErrInvalidFilters = fmt.Errorf("invalid filters")
)

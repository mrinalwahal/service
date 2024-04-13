package db

import "fmt"

var (
	ErrInvalidOptions = fmt.Errorf("invalid options")
	ErrInvalidUserID  = fmt.Errorf("invalid user_id")
	ErrInvalidTitle   = fmt.Errorf("invalid title")
	ErrInvalidFilters = fmt.Errorf("invalid filters")
)

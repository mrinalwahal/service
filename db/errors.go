package db

import "fmt"

var (
	ErrInvalidOptions  = fmt.Errorf("invalid options")
	ErrInvalidRecordID = fmt.Errorf("invalid record id")
	ErrInvalidUserID   = fmt.Errorf("invalid user id")
	ErrInvalidTitle    = fmt.Errorf("invalid title")
	ErrInvalidFilters  = fmt.Errorf("invalid filters")
	ErrNoRowsAffected  = fmt.Errorf("no rows affected")
)

package db

import "fmt"

var (
	ErrInvalidUserID = fmt.Errorf("invalid UserID")
	ErrEmptyTitle    = fmt.Errorf("empty title")
)

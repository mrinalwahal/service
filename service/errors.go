package service

import "fmt"

var (
	ErrInvalidOptions           = fmt.Errorf("invalid options")
	ErrOptionsNotFound          = fmt.Errorf("options not found")
	ErrRequesterDetailsNotFound = fmt.Errorf("requester details not found")
)

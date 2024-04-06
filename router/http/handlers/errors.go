package handlers

import "fmt"

var ErrInvalidRecordID = fmt.Errorf("invalid record id")
var ErrRecordNotFound = fmt.Errorf("record not found")

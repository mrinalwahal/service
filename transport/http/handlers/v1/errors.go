package v1

import "fmt"

var ErrInvalidRecordID = fmt.Errorf("invalid record id")
var ErrRecordNotFound = fmt.Errorf("record not found")
var ErrInvalidRequestOptions = fmt.Errorf("invalid request options")

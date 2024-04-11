package commons

import "github.com/google/uuid"

// Requester is the structure that holds the information of the user who sent the request.
type Requester struct {
	ID uuid.UUID
}

package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrinalwahal/service/pkg/middleware"
)

// Default HTTP Response structure.
// This structure implements the `error` interface.
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Err     error       `json:"error,omitempty"`
	Status  int         `json:"-"`
}

// Error returns the error message.
//
// This method is required to implement the `error` interface.
func (r *Response) Error() string {
	if r.Err != nil {
		return r.Err.Error()
	}
	return r.Message
}

func (r Response) MarshalJSON() ([]byte, error) {
	var errorMsg string
	if r.Err != nil {
		errorMsg = r.Err.Error()
	}
	var structure = struct {
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message,omitempty"`
		Err     string      `json:"error,omitempty"`
	}{
		Data:    r.Data,
		Message: r.Message,
		Err:     errorMsg,
	}
	return json.Marshal(structure)
}

func handleErr(w http.ResponseWriter, err error) {

	// Run type assertion on the response to check if it is of type `response`.
	// If it is, then write the response as JSON.
	// If it is not, then wrap the error in a new `Response` structure with defaults.
	if response, ok := err.(*Response); ok {
		if err := write(w, response.Status, response); err != nil {
			log.Println("failed to write response:", err)
		}
		return
	}
	if err := write(w, http.StatusInternalServerError, &Response{
		Message: "Your broke something on our server :(",
		Err:     err,
	}); err != nil {
		log.Println("failed to write response:", err)
	}
}

// write writes the data to the supplied http response writer.
func write(w http.ResponseWriter, status int, response any) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

// decode decodes the request body into the supplied type.
func decode[T any](r *http.Request) (T, error) {
	defer r.Body.Close()
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
}

// Validate the JWT Claims.
func (c *JWTClaims) Validate() error {
	if c.UserID == uuid.Nil {
		return ErrInvalidUserID
	}
	return nil
}

// Requester is the structure that holds the information of the user who sent the request.
type Requester struct {
	ID uuid.UUID
}

func getClaims(ctx context.Context) (*JWTClaims, error) {

	// Load the claims from request context to pass them in the service method.
	claims, ok := ctx.Value(middleware.XJWTClaims).(JWTClaims)
	if !ok {
		return nil, ErrInvalidJWTClaims
	}
	// Validate the claims.
	if err := claims.Validate(); err != nil {
		return nil, err
	}

	return &claims, nil
}

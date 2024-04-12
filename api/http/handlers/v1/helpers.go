package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Default HTTP Response structure.
// This structure implements the `error` interface.
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Err     error       `json:"error,omitempty"`
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

func (r *Response) UnmarshalJSON(data []byte) error {
	var structure = struct {
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message,omitempty"`
		Err     string      `json:"error,omitempty"`
	}{}
	if err := json.Unmarshal(data, &structure); err != nil {
		return err
	}
	r.Data = structure.Data
	r.Message = structure.Message
	if structure.Err != "" {
		r.Err = fmt.Errorf(structure.Err)
	}
	return nil
}

// write writes the data to the supplied http response writer.
func write(w http.ResponseWriter, status int, response any) error {
	w.WriteHeader(status)
	return encode(w, response)
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

// encode encodes the supplied data into the response writer.
func encode(w http.ResponseWriter, data any) error {
	return json.NewEncoder(w).Encode(data)
}

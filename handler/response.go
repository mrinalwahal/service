package handler

import (
	"encoding/json"
	"net/http"
)

// Default HTTP response structure.
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
	return r.Message
}

// WriteJSON writes the data to the supplied http response writer.
func WriteJSON(w http.ResponseWriter, status int, response any) error {

	// Set the content type to JSON.
	// w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

// // WriteString writes the data to the supplied http response writer.
// func WriteString(w http.ResponseWriter, status int, response string) error {
// 	w.WriteHeader(status)
// 	_, err := w.Write([]byte(response))
// 	return err
// }

// // WriteError writes the error to the supplied http response writer.
// func WriteError(w http.ResponseWriter, status int, response *Response) error {
// 	return WriteJSON(w, status, response)
// }

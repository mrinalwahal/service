package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Default API response structure.
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

// Convert the response to logger compatible.
func (r *Response) LogValue() slog.Value {

	//	Marshal the response to string.
	data, err := json.Marshal(r)
	if err != nil {
		return slog.StringValue("failed to marshal response")
	}

	return slog.StringValue(string(data))
}

// WriteJSON writes the data to the supplied http response writer.
func WriteJSON(w http.ResponseWriter, status int, response *Response) error {

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

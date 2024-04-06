package main

import (
	"encoding/json"
	"net/http"

	"github.com/mrinalwahal/service/handler"
)

//
// Utility functions.
//

func handle(handlerFunc func(*http.Request) error) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := handlerFunc(req); err != nil {

			// Run type assertion on the response to check if it is of type `Response`.
			// If it is, then write the response as JSON.
			// If it is not, then wrap the error in a new `Response` structure with defaults.
			if response, ok := err.(*handler.Response); ok {
				WriteJSON(w, response.Status, response)
				return
			}
			WriteJSON(w, http.StatusInternalServerError, &handler.Response{
				Message: "Your broke something on our server :(",
				Err:     err,
			})
			return
		}
	}
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

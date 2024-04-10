package v1

import (
	"net/http"
)

// Handler interface declares the signature of an HTTP request handler.
type Handler interface {

	// ServeHTTP is the method that consumes the incoming HTTP request.
	// Implementing this method in our custom handlers also implements the native `http.Handler` interface which is very good for compatibility with Go's inbuilt HTTP router.
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// Authorize is the method that runs custom validation on the incoming HTTP request.
	// This method is called before the `Process` method finally processes the request.
	//
	// Example uses:
	// - Check if the requester is authorized to make the request or not. You can implement access control over here.
	// - Check if the request body is valid or not.
	// - Check if the request headers are valid or not.
	// Authorize(context.Context) error

	// Validate is the method that runs custom validation on the incoming HTTP request.
	// This method is called before the `Process` method finally processes the request.
	//
	// Example uses:
	// - Check if the requester is authorized to make the request or not. You can implement access control over here.
	// - Check if the request body is valid or not.
	// - Check if the request headers are valid or not.
	// validate(context.Context, *CreateOptions) error

	// Process is the method that processes the incoming HTTP request to it's completion.
	// This method applies the main business logic of the handler.
	// And is called after the `Validate` method has successfully validated the request.
	// Process(context.Context) error
}

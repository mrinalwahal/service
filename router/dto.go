package router

import "encoding/json"

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

// CreateOptions represents the options for creating a todo.
type CreateOptions struct {

	//	Title of the todo.
	Title string `json:"title" validate:"required"`
}

// ListOptions represents the options for listing todos.
type ListOptions struct {

	//	Number of records to skip.
	Skip int `query:"skip" validate:"gte=0"`

	//	Number of records to return.
	Limit int `query:"limit" validate:"gte=0,lte=100"`

	//	Order by field.
	OrderBy string `query:"orderBy" validate:"oneof=created_at updated_at title"`

	//	Order by direction.
	OrderDirection string `query:"orderDirection" validate:"oneof=asc desc"`

	//	Title of the todo.
	Title string `query:"name"`
}

// UpdateOptions represents the options for updating a todo.
type UpdateOptions struct {

	//	Title of the todo.
	Title string `json:"title"`
}

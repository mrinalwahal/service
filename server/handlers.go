package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mrinalwahal/service/service"
)

// Create Handler.
func (s *HTTPServer) create(c echo.Context) error {

	//	Unmarshal the incoming payload.
	var payload CreateOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todo, err := svc.Create(ctx, &service.CreateOptions{
		Title: payload.Title,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to create the todo.",
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo created successfully.",
		Data:    todo,
	})
}

// Get Handler.
func (s *HTTPServer) get(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todo, err := svc.Get(ctx, uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo fetched successfully.",
		Data:    todo,
	})
}

// List Handler.
func (s *HTTPServer) list(c echo.Context) error {

	//	Unmarshal the incoming payload
	var payload ListOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todos, err := svc.List(ctx, &service.ListOptions{
		Skip:           payload.Skip,
		Limit:          payload.Limit,
		Title:          payload.Title,
		OrderBy:        payload.OrderBy,
		OrderDirection: payload.OrderDirection,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to list the todos.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todos fetched successfully.",
		Data:    todos,
	})
}

// Update Handler.
func (s *HTTPServer) update(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Unmarshal the incoming payload.
	var payload UpdateOptions
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid payload.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	todo, err := svc.Update(ctx, uuid, &service.UpdateOptions{
		Title: payload.Title,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo updated successfully.",
		Data:    todo,
	})
}

// Delete Handler.
func (s *HTTPServer) delete(c echo.Context) error {

	//	Get the object ID from the URL.
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID.")
	}

	//	Initialize a default context.
	ctx := context.Background()

	//	Get the service.
	svc, err := s.getService()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Error: "Failed to either connect to the database or prepare the service.",
		})
	}

	//	Call the service function to execute the business logic.
	err = svc.Delete(ctx, uuid)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create the todo.")
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "Todo deleted successfully.",
	})
}

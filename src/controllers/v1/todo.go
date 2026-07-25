// Package v1 holds versioned HTTP controllers. Controllers stay thin:
// decode, call the service, respond. No business logic lives here.
package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/furkancmn57/go-architecture-template/src/common"
	"github.com/furkancmn57/go-architecture-template/src/common/apperr"
	"github.com/furkancmn57/go-architecture-template/src/constants"
	"github.com/furkancmn57/go-architecture-template/src/models/requests"
	_ "github.com/furkancmn57/go-architecture-template/src/models/responses"
	todoservice "github.com/furkancmn57/go-architecture-template/src/services/todo"
)

// TodoController holds thin HTTP handlers for the todo resource.
type TodoController struct {
	service *todoservice.Service
}

// NewTodoController wires a TodoController to its service.
func NewTodoController(service *todoservice.Service) *TodoController {
	return &TodoController{service: service}
}

// Register mounts this controller's routes under the API group.
func (ctrl *TodoController) Register(api fiber.Router) {
	todos := api.Group("/todos")
	todos.Post("/", ctrl.Create)
	todos.Get("/", ctrl.List)
	todos.Get("/:id", ctrl.ByID)
	todos.Put("/:id", ctrl.Update)
	todos.Post("/:id/complete", ctrl.Complete)
	todos.Delete("/:id", ctrl.Delete)
}

// Create godoc
// @Summary      Create a todo
// @Description  Creates a new todo item
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        request body requests.CreateTodoRequest true "Create todo request"
// @Success      201 {object} responses.TodoResponse
// @Failure      422 {object} apperr.Error
// @Router       /todos [post]
func (ctrl *TodoController) Create(c *fiber.Ctx) error {
	var req requests.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.WriteHTTP(c, apperr.BadRequest(constants.InvalidRequestBody, "invalid request body"))
	}

	resp, appErr := ctrl.service.Create(c.Context(), req)
	if appErr != nil {
		return apperr.WriteHTTP(c, appErr)
	}

	return common.WriteJSON(c, fiber.StatusCreated, resp)
}

// List godoc
// @Summary      List todos
// @Tags         todos
// @Produce      json
// @Success      200 {array} responses.TodoResponse
// @Router       /todos [get]
func (ctrl *TodoController) List(c *fiber.Ctx) error {
	todos, appErr := ctrl.service.Todos(c.Context())
	if appErr != nil {
		return apperr.WriteHTTP(c, appErr)
	}
	return common.WriteJSON(c, fiber.StatusOK, todos)
}

// ByID godoc
// @Summary      Fetch a todo by id
// @Tags         todos
// @Produce      json
// @Param        id path string true "Todo ID"
// @Success      200 {object} responses.TodoResponse
// @Failure      404 {object} apperr.Error
// @Router       /todos/{id} [get]
func (ctrl *TodoController) ByID(c *fiber.Ctx) error {
	resp, appErr := ctrl.service.TodoById(c.Context(), c.Params("id"))
	if appErr != nil {
		return apperr.WriteHTTP(c, appErr)
	}
	return common.WriteJSON(c, fiber.StatusOK, resp)
}

// Update godoc
// @Summary      Update a todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id path string true "Todo ID"
// @Param        request body requests.UpdateTodoRequest true "Update todo request"
// @Success      200 {object} responses.TodoResponse
// @Failure      404 {object} apperr.Error
// @Router       /todos/{id} [put]
func (ctrl *TodoController) Update(c *fiber.Ctx) error {
	var req requests.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.WriteHTTP(c, apperr.BadRequest(constants.InvalidRequestBody, "invalid request body"))
	}

	resp, appErr := ctrl.service.Update(c.Context(), c.Params("id"), req)
	if appErr != nil {
		return apperr.WriteHTTP(c, appErr)
	}
	return common.WriteJSON(c, fiber.StatusOK, resp)
}

// Complete godoc
// @Summary      Mark a todo as completed
// @Tags         todos
// @Produce      json
// @Param        id path string true "Todo ID"
// @Success      200 {object} responses.TodoResponse
// @Failure      404 {object} apperr.Error
// @Router       /todos/{id}/complete [post]
func (ctrl *TodoController) Complete(c *fiber.Ctx) error {
	resp, appErr := ctrl.service.Complete(c.Context(), c.Params("id"))
	if appErr != nil {
		return apperr.WriteHTTP(c, appErr)
	}
	return common.WriteJSON(c, fiber.StatusOK, resp)
}

// Delete godoc
// @Summary      Delete a todo
// @Tags         todos
// @Param        id path string true "Todo ID"
// @Success      204
// @Failure      404 {object} apperr.Error
// @Router       /todos/{id} [delete]
func (ctrl *TodoController) Delete(c *fiber.Ctx) error {
	if appErr := ctrl.service.Delete(c.Context(), c.Params("id")); appErr != nil {
		return apperr.WriteHTTP(c, appErr)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

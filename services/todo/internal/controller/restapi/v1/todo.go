package v1

import (
	"net/http"

	"github.com/evrone/todo-svc/internal/controller/restapi/v1/request"
	"github.com/evrone/todo-svc/internal/entity"
	"github.com/gofiber/fiber/v2"
)

// @Summary     Create todo
// @Description Create a new todo item
// @ID          create-todo
// @Tags        todo
// @Accept      json
// @Produce     json
// @Param       request body request.CreateTodo true "Todo to create"
// @Success     201 {object} entity.Todo
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /todo [post]
func (r *V1) create(ctx *fiber.Ctx) error {
	var body request.CreateTodo

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "restapi - v1 - todo - create")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - todo - create")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	t, err := r.todo.Create(
		ctx.UserContext(),
		entity.Todo{
			Title:       body.Title,
			Description: body.Description,
			Priority:    body.Priority,
			DueDate:     body.DueDate,
		},
	)
	if err != nil {
		r.l.Error(err, "restapi - v1 - todo - create")

		return errorResponse(ctx, http.StatusInternalServerError, "todo service problems")
	}

	return ctx.Status(http.StatusCreated).JSON(t)
}

// @Summary     Get todo by ID
// @Description Get a single todo item by its ID
// @ID          get-todo
// @Tags        todo
// @Produce     json
// @Param       id path int true "Todo ID"
// @Success     200 {object} entity.Todo
// @Failure     400 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /todo/{id} [get]
func (r *V1) getByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id param")
	}

	t, err := r.todo.GetByID(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "restapi - v1 - todo - getByID")

		return errorResponse(ctx, http.StatusInternalServerError, "todo service problems")
	}

	return ctx.Status(http.StatusOK).JSON(t)
}

// @Summary     List todos
// @Description Get all todo items
// @ID          list-todos
// @Tags        todo
// @Produce     json
// @Success     200 {array} entity.Todo
// @Failure     500 {object} response.Error
// @Router      /todo [get]
func (r *V1) list(ctx *fiber.Ctx) error {
	todos, err := r.todo.List(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "restapi - v1 - todo - list")

		return errorResponse(ctx, http.StatusInternalServerError, "todo service problems")
	}

	return ctx.Status(http.StatusOK).JSON(todos)
}

// @Summary     Update todo
// @Description Update an existing todo item
// @ID          update-todo
// @Tags        todo
// @Accept      json
// @Produce     json
// @Param       id      path int                true "Todo ID"
// @Param       request body request.UpdateTodo true "Fields to update"
// @Success     200 {object} entity.Todo
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /todo/{id} [put]
func (r *V1) update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id param")
	}

	var body request.UpdateTodo

	if err = ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "restapi - v1 - todo - update")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err = r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - todo - update")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	t, err := r.todo.Update(
		ctx.UserContext(),
		id,
		entity.Todo{
			Title:       body.Title,
			Description: body.Description,
			Status:      body.Status,
			Priority:    body.Priority,
			DueDate:     body.DueDate,
		},
	)
	if err != nil {
		r.l.Error(err, "restapi - v1 - todo - update")

		return errorResponse(ctx, http.StatusInternalServerError, "todo service problems")
	}

	return ctx.Status(http.StatusOK).JSON(t)
}

// @Summary     Delete todo
// @Description Delete a todo item by its ID
// @ID          delete-todo
// @Tags        todo
// @Param       id path int true "Todo ID"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /todo/{id} [delete]
func (r *V1) delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id param")
	}

	if err = r.todo.Delete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "restapi - v1 - todo - delete")

		return errorResponse(ctx, http.StatusInternalServerError, "todo service problems")
	}

	return ctx.SendStatus(http.StatusNoContent)
}

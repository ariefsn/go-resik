package api

import (
	"errors"
	"net/http"

	"github.com/ariefsn/go-resik/common"
	"github.com/ariefsn/go-resik/domain"
	"github.com/ariefsn/go-resik/helper"
	"github.com/ariefsn/go-resik/logger"
	"github.com/gofiber/fiber/v2"
)

// TodoApi  represent the httphandler for todo
type TodoApi struct {
	todoSvc domain.TodoService
}

func NewTodoApi(todoSvc domain.TodoService) *fiber.App {
	api := &TodoApi{
		todoSvc: todoSvc,
	}

	app := fiber.New()

	app.Post("/", api.Create).Name("todoCreate")
	app.Get("/", api.Get).Name("todoGet")
	app.Get("/:id", api.GetByID).Name("todoGetById")
	app.Put("/:id", api.Update)
	app.Patch("/:id", api.UpdateStatus)
	app.Delete("/:id", api.Delete)

	return app
}

func (a *TodoApi) Create(c *fiber.Ctx) error {
	payload := domain.TodoDto{}

	if err := c.BodyParser(&payload); err != nil {
		logger.Error(err)
		return c.Status(http.StatusBadRequest).JSON(helper.JsonError(err))
	}

	res, err := a.todoSvc.Create(c.Context(), &payload)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.JsonError(err))
	}

	return c.Status(http.StatusOK).JSON(helper.JsonSuccess(res))
}

func (a *TodoApi) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		err := errors.New("id is required")
		logger.Error(err)

		return c.Status(http.StatusBadRequest).JSON(helper.JsonError(err))
	}

	res, err := a.todoSvc.GetByID(c.Context(), id)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.JsonError(err))
	}

	return c.Status(http.StatusOK).JSON(helper.JsonSuccess(res))
}

func (a *TodoApi) Get(c *fiber.Ctx) error {
	skip := c.QueryInt("skip", 0)
	limit := c.QueryInt("limit", 10)

	title := c.Query("title")
	description := c.Query("description")

	filter := common.M{}

	if title != "" {
		filter["title"] = title
	}

	if description != "" {
		filter["description"] = description
	}

	res, total, err := a.todoSvc.Get(c.Context(), filter, int64(skip), int64(limit))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.JsonError(err))
	}

	return c.Status(http.StatusOK).JSON(helper.JsonSuccess(common.M{
		"items": res,
		"total": total,
	}))
}

func (a *TodoApi) Update(c *fiber.Ctx) error {
	payload := domain.TodoDto{}

	id := c.Params("id")

	if id == "" {
		err := errors.New("id is required")
		logger.Error(err)

		return c.Status(http.StatusBadRequest).JSON(helper.JsonError(err))
	}

	if err := c.BodyParser(&payload); err != nil {
		logger.Error(err)
		return c.Status(http.StatusBadRequest).JSON(helper.JsonError(err))
	}

	res, err := a.todoSvc.Update(c.Context(), id, &payload)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.JsonError(err))
	}

	return c.Status(http.StatusOK).JSON(helper.JsonSuccess(res))
}

func (a *TodoApi) UpdateStatus(c *fiber.Ctx) error {
	payload := struct {
		IsCompleted bool `json:"isCompleted"`
	}{}

	id := c.Params("id")

	if id == "" {
		err := errors.New("id is required")
		logger.Error(err)

		return c.Status(http.StatusBadRequest).JSON(helper.JsonError(err))
	}

	if err := c.BodyParser(&payload); err != nil {
		logger.Error(err)
		return c.Status(http.StatusBadRequest).JSON(helper.JsonError(err))
	}

	res, err := a.todoSvc.UpdateStatus(c.Context(), id, payload.IsCompleted)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.JsonError(err))
	}

	return c.Status(http.StatusOK).JSON(helper.JsonSuccess(res))
}

func (a *TodoApi) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		err := errors.New("id is required")
		logger.Error(err)

		return c.Status(http.StatusBadRequest).JSON(helper.JsonError(err))
	}

	err := a.todoSvc.Delete(c.Context(), id)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.JsonError(err))
	}

	return c.Status(http.StatusOK).JSON(helper.JsonSuccess(id))
}

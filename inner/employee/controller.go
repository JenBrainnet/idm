package employee

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
}

type Svc interface {
	Create(request CreateRequest) (int64, error)
	FindById(request IdRequest) (Response, error)
	FindAll() ([]Response, error)
	FindAllByIds(request IdsRequest) ([]Response, error)
	DeleteById(request IdRequest) error
	DeleteAllByIds(request IdsRequest) error
}

func NewController(server *web.Server, employeeService Svc) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/employees", c.CreateEmployee)
	c.server.GroupApiV1.Get("/employees/:id", c.FindById)
	c.server.GroupApiV1.Get("/employees", c.FindAll)
	c.server.GroupApiV1.Post("/employees/ids", c.FindAllByIds)
	c.server.GroupApiV1.Delete("/employees/:id", c.DeleteById)
	c.server.GroupApiV1.Delete("/employees", c.DeleteAllByIds)
}

func (c *Controller) CreateEmployee(ctx *fiber.Ctx) error {
	var request CreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	var newEmployeeId, err = c.employeeService.Create(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}

	if err = common.OkResponse(ctx, newEmployeeId); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
	}
	return nil
}

func (c *Controller) FindById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	response, err := c.employeeService.FindById(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindAll(ctx *fiber.Ctx) error {
	responses, err := c.employeeService.FindAll()
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse(ctx, responses)
}

func (c *Controller) FindAllByIds(ctx *fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	responses, err := c.employeeService.FindAllByIds(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse(ctx, responses)
}

func (c *Controller) DeleteById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	err = c.employeeService.DeleteById(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse[any](ctx, nil)
}

func (c *Controller) DeleteAllByIds(ctx *fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	err := c.employeeService.DeleteAllByIds(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse[any](ctx, nil)
}

func resolveHttpStatusCode(err error) int {
	switch {
	case errors.As(err, &common.RequestValidationError{}) || errors.As(err, &common.AlreadyExistsError{}):
		return fiber.StatusBadRequest
	case errors.As(err, &common.NotFoundError{}):
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}

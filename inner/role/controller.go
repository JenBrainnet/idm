package role

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
)

type Controller struct {
	server      *web.Server
	roleService Svc
}

type Svc interface {
	Create(request CreateRequest) (int64, error)
	FindById(request IdRequest) (Response, error)
	FindAll() ([]Response, error)
	FindAllByIds(request IdsRequest) ([]Response, error)
	DeleteById(request IdRequest) error
	DeleteAllByIds(request IdsRequest) error
}

func NewController(server *web.Server, roleService Svc) *Controller {
	return &Controller{
		server:      server,
		roleService: roleService,
	}
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupApiV1.Post("/roles", c.CreateRole)
	c.server.GroupApiV1.Get("/roles/:id", c.FindById)
	c.server.GroupApiV1.Get("/roles", c.FindAll)
	c.server.GroupApiV1.Post("/roles/ids", c.FindAllByIds)
	c.server.GroupApiV1.Delete("/roles/:id", c.DeleteById)
	c.server.GroupApiV1.Delete("/roles", c.DeleteAllByIds)
}

func (c *Controller) CreateRole(ctx fiber.Ctx) error {
	var request CreateRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	var newRoleId, err = c.roleService.Create(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}

	if err = common.OkResponse(ctx, newRoleId); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created role id")
	}
	return nil
}

func (c *Controller) FindById(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	response, err := c.roleService.FindById(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindAll(ctx fiber.Ctx) error {
	responses, err := c.roleService.FindAll()
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse(ctx, responses)
}

func (c *Controller) FindAllByIds(ctx fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	responses, err := c.roleService.FindAllByIds(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse(ctx, responses)
}

func (c *Controller) DeleteById(ctx fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	err = c.roleService.DeleteById(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse[any](ctx, nil)
}

func (c *Controller) DeleteAllByIds(ctx fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	err := c.roleService.DeleteAllByIds(request)
	if err != nil {
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	return common.OkResponse[any](ctx, nil)
}

func resolveHttpStatusCode(err error) int {
	switch {
	case errors.As(err, &common.RequestValidationError{}):
		return fiber.StatusBadRequest
	case errors.As(err, &common.NotFoundError{}):
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}

package role

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
)

type Controller struct {
	server      *web.Server
	roleService Svc
	logger      *common.Logger
}

type Svc interface {
	Create(request CreateRequest) (int64, error)
	FindById(request IdRequest) (Response, error)
	FindAll() ([]Response, error)
	FindAllByIds(request IdsRequest) ([]Response, error)
	DeleteById(request IdRequest) error
	DeleteAllByIds(request IdsRequest) error
}

func NewController(server *web.Server, roleService Svc, logger *common.Logger) *Controller {
	return &Controller{
		server:      server,
		roleService: roleService,
		logger:      logger,
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

func (c *Controller) CreateRole(ctx *fiber.Ctx) error {
	var request CreateRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("create role: failed to parse request body", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("create role: received request", zap.Any("request", request))
	var newRoleId, err = c.roleService.Create(request)
	if err != nil {
		c.logger.Error("create role: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("create role: success", zap.Int64("id", newRoleId))
	if err = common.OkResponse(ctx, newRoleId); err != nil {
		c.logger.Error("create role: failed to send response", zap.Int64("id", newRoleId), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created role id")
	}
	return nil
}

func (c *Controller) FindById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	c.logger.Debug("find role by id: received id", zap.String("id", idStr))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.logger.Error("find role by id: invalid id parameter", zap.String("id", idStr), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	response, err := c.roleService.FindById(request)
	if err != nil {
		c.logger.Error("find role by id: service error", zap.Int64("id", id), zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("find role by id: success", zap.Int64("id", id))
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindAll(ctx *fiber.Ctx) error {
	c.logger.Debug("find all roles: received request")
	responses, err := c.roleService.FindAll()
	if err != nil {
		c.logger.Error("find all roles: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("find all roles: success", zap.Int("count", len(responses)))
	return common.OkResponse(ctx, responses)
}

func (c *Controller) FindAllByIds(ctx *fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("find roles by ids: failed to parse request", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("find roles by ids: received request", zap.Any("request", request))
	responses, err := c.roleService.FindAllByIds(request)
	if err != nil {
		c.logger.Error("find roles by ids: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("find roles by ids: success", zap.Int("count", len(responses)))
	return common.OkResponse(ctx, responses)
}

func (c *Controller) DeleteById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	c.logger.Debug("delete role by id: received id", zap.String("id", idStr))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.logger.Error("delete role by id: invalid id parameter", zap.String("id", idStr), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	err = c.roleService.DeleteById(request)
	if err != nil {
		c.logger.Error("delete role by id: service error", zap.Int64("id", id), zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("delete role by id: success", zap.Int64("id", id))
	return common.OkResponse[any](ctx, nil)
}

func (c *Controller) DeleteAllByIds(ctx *fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("delete roles by ids: failed to parse request", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("delete roles by ids: received request", zap.Any("request", request))
	err := c.roleService.DeleteAllByIds(request)
	if err != nil {
		c.logger.Error("delete roles by ids: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("delete roles by ids: success", zap.Int("count", len(request.Ids)))
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

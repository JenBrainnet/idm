package employee

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/web"
	"strconv"
)

type Controller struct {
	server          *web.Server
	employeeService Svc
	logger          *common.Logger
}

type Svc interface {
	Create(request CreateRequest) (int64, error)
	FindById(request IdRequest) (Response, error)
	FindAll() ([]Response, error)
	FindAllByIds(request IdsRequest) ([]Response, error)
	DeleteById(request IdRequest) error
	DeleteAllByIds(request IdsRequest) error
}

func NewController(server *web.Server, employeeService Svc, logger *common.Logger) *Controller {
	return &Controller{
		server:          server,
		employeeService: employeeService,
		logger:          logger,
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
		c.logger.Error("create employee: failed to parse request body", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	c.logger.Debug("create employee: received request", zap.Any("request", request))
	var newEmployeeId, err = c.employeeService.Create(request)
	if err != nil {
		c.logger.Error("create employee: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}

	c.logger.Debug("create employee: success", zap.Int64("id", newEmployeeId))
	if err = common.OkResponse(ctx, newEmployeeId); err != nil {
		c.logger.Error("create employee: failed to send response")
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning created employee id")
	}
	return nil
}

func (c *Controller) FindById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	c.logger.Debug("find employee by id: received id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.logger.Error("find employee by id: invalid id parameter", zap.String("id", idStr), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	response, err := c.employeeService.FindById(request)
	if err != nil {
		c.logger.Error("find employee by id: service error", zap.Int64("id", id), zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("find employee by id: success", zap.Int64("id", id))
	return common.OkResponse(ctx, response)
}

func (c *Controller) FindAll(ctx *fiber.Ctx) error {
	c.logger.Debug("find all employees: received request")
	responses, err := c.employeeService.FindAll()
	if err != nil {
		c.logger.Error("find all employees: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("find all employees: success", zap.Int("count", len(responses)))
	return common.OkResponse(ctx, responses)
}

func (c *Controller) FindAllByIds(ctx *fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("find all employees by ids: body parse error", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("find all employees by ids: received request", zap.Any("request", request))
	responses, err := c.employeeService.FindAllByIds(request)
	if err != nil {
		c.logger.Error("find all employees by ids: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("find all employees by ids: success", zap.Int("count", len(responses)))
	return common.OkResponse(ctx, responses)
}

func (c *Controller) DeleteById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	c.logger.Debug("delete employee by id: received id", zap.String("id", idStr))
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.logger.Error("delete employee by id: invalid id parameter", zap.String("id", idStr), zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, "invalid id parameter")
	}
	request := IdRequest{Id: id}
	err = c.employeeService.DeleteById(request)
	if err != nil {
		c.logger.Error("delete employee by id: service error", zap.Int64("id", id), zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("delete employee by id: success", zap.Int64("id", id))
	return common.OkResponse[any](ctx, nil)
}

func (c *Controller) DeleteAllByIds(ctx *fiber.Ctx) error {
	var request IdsRequest
	if err := ctx.BodyParser(&request); err != nil {
		c.logger.Error("delete all employees by ids: body parse error", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusBadRequest, err.Error())
	}
	c.logger.Debug("delete all employees by ids: received request", zap.Any("request", request))
	err := c.employeeService.DeleteAllByIds(request)
	if err != nil {
		c.logger.Error("delete all employees by ids: service error", zap.Error(err))
		return common.ErrResponse(ctx, resolveHttpStatusCode(err), err.Error())
	}
	c.logger.Debug("delete all employees by ids: success", zap.Int("count", len(request.Ids)))
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

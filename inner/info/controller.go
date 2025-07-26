package info

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/web"
	"time"
)

type Controller struct {
	server *web.Server
	cfg    common.Config
	db     *sqlx.DB
	logger *common.Logger
}

func NewController(server *web.Server, cfg common.Config, db *sqlx.DB, logger *common.Logger) *Controller {
	return &Controller{
		server: server,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

type InfoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Controller) RegisterRoutes() {
	c.server.GroupInternal.Get("/info", c.GetInfo)
	c.server.GroupInternal.Get("/health", c.GetHealth)
}

func (c *Controller) GetInfo(ctx *fiber.Ctx) error {
	c.logger.Debug("get info: received request")
	err := ctx.Status(fiber.StatusOK).JSON(&InfoResponse{
		Name:    c.cfg.AppName,
		Version: c.cfg.AppVersion,
	})
	if err != nil {
		c.logger.Error("get info: failed to send response", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning info")
	}
	c.logger.Debug("get info: success", zap.String("name", c.cfg.AppName), zap.String("version", c.cfg.AppVersion))
	return nil
}

func (c *Controller) GetHealth(ctx *fiber.Ctx) error {
	c.logger.Debug("get health: received request")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := c.db.PingContext(ctxTimeout); err != nil {
		c.logger.Error("get health: database unreachable", zap.Error(err))
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "DB not reachable")
	}
	c.logger.Debug("get health: success")
	return ctx.Status(fiber.StatusOK).SendString("OK")
}

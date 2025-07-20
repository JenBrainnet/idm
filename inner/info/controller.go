package info

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
	"idm/inner/web"
	"time"
)

type Controller struct {
	server *web.Server
	cfg    common.Config
	db     *sqlx.DB
}

func NewController(server *web.Server, cfg common.Config, db *sqlx.DB) *Controller {
	return &Controller{
		server: server,
		cfg:    cfg,
		db:     db,
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
	err := ctx.Status(fiber.StatusOK).JSON(&InfoResponse{
		Name:    c.cfg.AppName,
		Version: c.cfg.AppVersion,
	})
	if err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "error returning info")
	}
	return nil
}

func (c *Controller) GetHealth(ctx *fiber.Ctx) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := c.db.PingContext(ctxTimeout); err != nil {
		return common.ErrResponse(ctx, fiber.StatusInternalServerError, "DB not reachable")
	}
	return ctx.Status(fiber.StatusOK).SendString("OK")
}

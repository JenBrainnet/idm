package web

import "github.com/gofiber/fiber/v3"

type Server struct {
	App        *fiber.App
	GroupApiV1 fiber.Router
}

func NewServer() *Server {
	app := fiber.New()
	groupApi := app.Group("/api")
	groupApiV1 := groupApi.Group("/v1")

	return &Server{
		App:        app,
		GroupApiV1: groupApiV1,
	}
}

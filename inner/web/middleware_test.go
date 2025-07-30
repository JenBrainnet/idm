package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverMiddleware(t *testing.T) {
	app := fiber.New()
	RegisterMiddleware(app)

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "test panic")
}

func TestRequestIDMiddleware(t *testing.T) {
	app := fiber.New()
	RegisterMiddleware(app)

	app.Get("/id", func(c *fiber.Ctx) error {
		id := c.GetRespHeader(fiber.HeaderXRequestID)
		return c.SendString(id)
	})

	req := httptest.NewRequest(http.MethodGet, "/id", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	requestID := resp.Header.Get(fiber.HeaderXRequestID)
	assert.NotEmpty(t, requestID)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, requestID, string(body))
}

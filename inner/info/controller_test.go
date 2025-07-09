package info_test

import (
	"encoding/json"
	"errors"
	"github.com/78bits/go-sqlmock-sqlx"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"idm/inner/info"
	"idm/inner/web"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInfo(t *testing.T) {
	a := assert.New(t)

	t.Run("should return info response", func(t *testing.T) {
		cfg := common.Config{
			AppName:    "idm-service",
			AppVersion: "1.0.0",
		}
		app := fiber.New()
		server := &web.Server{
			App:           app,
			GroupInternal: app.Group("/internal"),
		}
		controller := info.NewController(server, cfg, nil)
		controller.RegisterRoutes()

		req := httptest.NewRequest(http.MethodGet, "/internal/info", nil)
		resp, err := app.Test(req)
		a.NoError(err)
		a.Equal(http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response info.InfoResponse
		err = json.Unmarshal(body, &response)
		a.NoError(err)
		a.Equal("idm-service", response.Name)
		a.Equal("1.0.0", response.Version)
	})
}

func TestGetHealth(t *testing.T) {
	a := assert.New(t)

	t.Run("should return OK when DB is reachable", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		a.NoError(err)
		defer db.Close()
		mock.ExpectPing().WillReturnError(nil)
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		app := fiber.New()
		server := &web.Server{
			App:           app,
			GroupInternal: app.Group("/internal"),
		}
		cfg := common.Config{}
		controller := info.NewController(server, cfg, sqlxDB)
		controller.RegisterRoutes()

		req := httptest.NewRequest(http.MethodGet, "/internal/health", nil)
		resp, err := app.Test(req)
		a.NoError(err)
		a.Equal(http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		a.NoError(err)
		a.Equal("OK", string(body))
	})

	t.Run("should return internal server error when DB is unreachable", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		a.NoError(err)
		defer db.Close()
		mock.ExpectPing().WillReturnError(errors.New("connection refused"))
		sqlxDB := sqlx.NewDb(db, "sqlmock")
		app := fiber.New()
		server := &web.Server{
			App:           app,
			GroupInternal: app.Group("/internal"),
		}
		cfg := common.Config{}
		controller := info.NewController(server, cfg, sqlxDB)
		controller.RegisterRoutes()

		req := httptest.NewRequest(http.MethodGet, "/internal/health", nil)
		resp, err := app.Test(req)
		a.NoError(err)
		a.Equal(http.StatusInternalServerError, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		a.NoError(err)
		var errResponse common.Response[any]
		err = json.Unmarshal(body, &errResponse)
		a.NoError(err)
		a.False(errResponse.Success)
		a.Equal("DB not reachable", errResponse.Message)
	})
}

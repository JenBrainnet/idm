package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg := common.GetConfig(".env")
	logger := common.NewLogger(cfg)
	defer func() {
		_ = logger.Sync()
	}()
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error closing db: %v", zap.Error(err))
		}
	}()
	server := build(cfg, db, logger)
	go func() {
		err := server.App.Listen(":8080")
		if err != nil {
			logger.Panic("http server error: %s", zap.Error(err))
		}
	}()
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go gracefulShutdown(server, wg, logger)
	wg.Wait()
	logger.Info("Graceful shutdown complete.")
}

func gracefulShutdown(server *web.Server, wg *sync.WaitGroup, logger *common.Logger) {
	defer wg.Done()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()
	<-ctx.Done()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown with error: %v\n", zap.Error(err))
	}
	logger.Info("Server exiting")
}

func build(cfg common.Config, db *sqlx.DB, logger *common.Logger) *web.Server {
	server := web.NewServer()
	employeeRepo := employee.NewRepository(db)
	roleRepo := role.NewRepository(db)
	vld := validator.New()
	employeeService := employee.NewService(employeeRepo, vld)
	roleService := role.NewService(roleRepo, vld)
	employeeController := employee.NewController(server, employeeService, logger)
	roleController := role.NewController(server, roleService, logger)
	employeeController.RegisterRoutes()
	roleController.RegisterRoutes()
	infoController := info.NewController(server, cfg, db, logger)
	infoController.RegisterRoutes()
	return server
}

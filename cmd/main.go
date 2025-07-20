package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
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
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("error closing db: %v", err)
		}
	}()
	server := build(db)
	go func() {
		err := server.App.Listen(":8080")
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	go gracefulShutdown(server, wg)
	wg.Wait()
	fmt.Println("Graceful shutdown complete.")
}

func gracefulShutdown(server *web.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()
	<-ctx.Done()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("Server forced to shutdown with error: %v\n", err)
	}
	fmt.Println("Server exiting")
}

func build(db *sqlx.DB) *web.Server {
	cfg := common.GetConfig(".env")
	server := web.NewServer()
	employeeRepo := employee.NewRepository(db)
	roleRepo := role.NewRepository(db)
	vld := validator.New()
	employeeService := employee.NewService(employeeRepo, vld)
	roleService := role.NewService(roleRepo, vld)
	employeeController := employee.NewController(server, employeeService)
	roleController := role.NewController(server, roleService)
	employeeController.RegisterRoutes()
	roleController.RegisterRoutes()
	infoController := info.NewController(server, cfg, db)
	infoController.RegisterRoutes()
	return server
}

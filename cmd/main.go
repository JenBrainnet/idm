package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
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
	err := server.App.Listen(":8080")
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
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

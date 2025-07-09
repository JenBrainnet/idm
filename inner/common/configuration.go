package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

// Config общая конфигурация всего приложения
type Config struct {
	DbDriverName string `validate:"required"`
	Dsn          string `validate:"required"`
	AppName      string `validate:"required"`
	AppVersion   string `validate:"required"`
}

// GetConfig получение конфигурации из .env файла или переменных окружения
func GetConfig(envFile string) Config {
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}
	var cfg = Config{
		DbDriverName: os.Getenv("DB_DRIVER_NAME"),
		Dsn:          os.Getenv("DB_DSN"),
		AppName:      os.Getenv("APP_NAME"),
		AppVersion:   os.Getenv("APP_VERSION"),
	}
	err = validator.New().Struct(cfg)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			panic(fmt.Sprintf("config validation error: %v", err))
		}
	}
	return cfg
}

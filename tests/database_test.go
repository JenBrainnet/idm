package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"idm/inner/database"
	"os"
	"testing"
)

const testDbDsn = "host=localhost port=5432 user=postgres password=postgres dbname=idm_test_db sslmode=disable"

func setTestEnv(driver, dsn string) {
	_ = os.Setenv("DB_DRIVER_NAME", driver)
	_ = os.Setenv("DB_DSN", dsn)
	_ = os.Setenv("APP_NAME", "idm")
	_ = os.Setenv("APP_VERSION", "0.0.0")
}

func unsetTestEnv() {
	_ = os.Unsetenv("DB_DRIVER_NAME")
	_ = os.Unsetenv("DB_DSN")
	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")
}

// В проекте нет .env  файла (должны получить конфигурацию из переменных окружения)
func TestConfigNoEnvFileUseEnvVars(t *testing.T) {
	a := assert.New(t)

	setTestEnv("postgres", testDbDsn)
	defer unsetTestEnv()

	cfg := common.GetConfig("nonexisted.env")

	a.Equal("postgres", cfg.DbDriverName)
	a.Equal(testDbDsn, cfg.Dsn)
	a.Equal("idm", cfg.AppName)
	a.Equal("0.0.0", cfg.AppVersion)
}

func TestConfigEmptyEnvFileAndNoEnvVars(t *testing.T) {
	a := assert.New(t)
	defer func() {
		r := recover()
		a.NotNil(r, "Expected panic due to missing required config")
		a.Contains(fmt.Sprintf("%v", r), "AppName") // можно проверить текст panic, если нужно
	}()

	tmp, _ := os.CreateTemp("", ".env")
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmp.Name()) }()

	unsetTestEnv()

	_ = common.GetConfig(tmp.Name())
}

// В проекте есть .env  файл и в нём нет нужных переменных, но в переменных окружения
// они есть (должны получить заполненную структуру  idm.inner.common.Config с данными
// из пременных окружения)
func TestConfigEmptyEnvFileUseEnvVars(t *testing.T) {
	a := assert.New(t)

	tmp, _ := os.CreateTemp("", ".env")
	_, _ = tmp.WriteString("FOO=bar")
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmp.Name()) }()

	setTestEnv("postgres", testDbDsn)
	defer unsetTestEnv()

	cfg := common.GetConfig(tmp.Name())

	a.Equal("postgres", cfg.DbDriverName)
	a.Equal(testDbDsn, cfg.Dsn)
}

// В проекте есть корректно заполненный .env файл, в переменных окружения нет конфликтующих
// с ним переменных  (должны получить структуру  idm.inner.common.Config, заполненную данными
// из .env файла)
func TestConfigNoEnvVarsUseEnvFile(t *testing.T) {
	a := assert.New(t)

	tmp, _ := os.CreateTemp("", ".env")
	_, _ = tmp.WriteString("DB_DRIVER_NAME=postgres\nDB_DSN=" + testDbDsn + "\nAPP_NAME=idm\nAPP_VERSION=0.0.0\n")
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmp.Name()) }()

	unsetTestEnv()

	cfg := common.GetConfig(tmp.Name())

	a.Equal("postgres", cfg.DbDriverName)
	a.Equal(testDbDsn, cfg.Dsn)
}

// В проекте есть .env  файл и в нём есть нужные переменные, но в переменных окружения они
// тоже есть (с другими значениями) - должны получить структуру  idm.inner.common.Config,
// заполненную данными. Используются значения из переменных окружения
func TestConfigEnvVarsAndEnvFileExistUseEnvVars(t *testing.T) {
	a := assert.New(t)

	tmp, _ := os.CreateTemp("", ".env")
	_, _ = tmp.WriteString("DB_DRIVER_NAME=file-driver-name\nDB_DSN=file-dsn")
	_ = tmp.Close()
	defer func() { _ = os.Remove(tmp.Name()) }()

	setTestEnv("env-vars-driver-name", "env-vars-dsn")
	defer unsetTestEnv()

	cfg := common.GetConfig(tmp.Name())

	a.Equal("env-vars-driver-name", cfg.DbDriverName)
	a.Equal("env-vars-dsn", cfg.Dsn)
}

// Приложение не может подключиться к базе данных с некорректным конфигом
func TestConnectDbWithInvalidPortFails(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic when connecting with invalid config, " +
				"but no panic occurred")
		}
	}()
	cfg := common.Config{
		DbDriverName: "postgres",
		Dsn:          "host=localhost port=0000 user=postgres password=postgres dbname=idm_test_db sslmode=disable",
		AppName:      "idm",
		AppVersion:   "0.0.0",
	}
	_ = database.ConnectDbWithCfg(cfg)
}

// Приложение может подключиться к базе данных с корректным конфигом.
func TestConnectDbWithValidConfigOk(t *testing.T) {
	a := assert.New(t)

	cfg := common.Config{
		DbDriverName: "postgres",
		Dsn:          testDbDsn,
		AppName:      "idm",
		AppVersion:   "0.0.0",
	}

	db := database.ConnectDbWithCfg(cfg)
	defer func() { _ = db.Close() }()

	a.NotNil(db)
	err := db.Ping()
	a.NoError(err)
}

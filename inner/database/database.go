package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"idm/inner/common"
	"time"
)

// ConnectDb получить конфиг и подключиться с ним к базе данных
func ConnectDb() *sqlx.DB {
	cfg := common.GetConfig(".env")
	return ConnectDbWithCfg(cfg)
}

// ConnectDbWithCfg подключиться к базе данных с переданным конфигом
func ConnectDbWithCfg(cfg common.Config) *sqlx.DB {
	var db = sqlx.MustConnect(cfg.DbDriverName, cfg.Dsn)
	// Настройки ниже конфигурируют пулл подключений к базе данных. Их названия стандартны для большинства библиотек.
	// Ознакомиться с их описанием можно на примере документации Hikari pool:
	// https://github.com/brettwooldridge/HikariCP?tab=readme-ov-file#gear-configuration-knobs-baby
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	return db
}

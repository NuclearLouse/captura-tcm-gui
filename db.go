package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type dB struct {
	Connect *gorm.DB
}

func newDB(cfg *config) (*dB, error) {
	db, err := gorm.Open(connectOptions(cfg))
	if err != nil {
		return nil, err
	}
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}
	return &dB{Connect: db}, nil
}

func connectOptions(cfg *config) (string, string) {
	var dialect, sslmode, connectOptions string
	switch cfg.ConnectDB.SslMode {
	case true:
		sslmode = "enable"
	case false:
		sslmode = "disable"
	}
	dialect = os.Getenv("DIALECT_DB")
	switch dialect {
	case "mysql":
		connectOptions = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=%s",
			cfg.ConnectDB.User,
			cfg.ConnectDB.Pass,
			cfg.ConnectDB.Database,
			cfg.ConnectDB.Host+":"+cfg.ConnectDB.Port)
	case "postgres":
		connectOptions = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.ConnectDB.Host,
			cfg.ConnectDB.Port,
			cfg.ConnectDB.User,
			cfg.ConnectDB.Database,
			cfg.ConnectDB.Pass,
			sslmode)
	case "sqlite3":
		connectOptions = "tcm.db"
	}
	return dialect, connectOptions
}

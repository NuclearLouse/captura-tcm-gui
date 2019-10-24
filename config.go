package main

import (
	"os"

	"github.com/go-ini/ini"
)

type config struct {
	ConnectDB `ini:"connectdb"`
}

type ConnectDB struct {
	Dialect      string `ini:"db_dialect"`
	Host         string `ini:"host"`
	Port         string `ini:"port"`
	Database     string `ini:"database"`
	SchemaPG     string `ini:"schema_pg"`
	User         string `ini:"user"`
	Pass         string `ini:"password"`
	SslMode      bool   `ini:"ssl_mode"`
	CreateTables bool   `ini:"create_tables"`
}

func loadFile(configFile string) (*ini.File, error) {
	f, err := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, configFile)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func readConfig(configFile string) (*config, error) {
	cfg := &config{}
	set, err := loadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err = set.MapTo(&cfg); err != nil {
		return nil, err
	}
	setEnvVars(cfg)
	return cfg, nil
}

func setEnvVars(cfg *config) {
	var schema, dialect string
	schema = cfg.ConnectDB.SchemaPG + "."
	if schema == "" {
		schema = "public."
	}
	dialect = cfg.ConnectDB.Dialect
	if dialect == "sqlite" {
		dialect = "sqlite3"
	}
	if dialect == "sqlite3" || dialect == "mysql" {
		schema = ""
	}
	os.Setenv("SCHEMA_PG", schema)
	os.Setenv("DIALECT_DB", dialect)
}

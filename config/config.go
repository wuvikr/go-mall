package config

import "time"

var (
	App      *appConfig
	Database *databaseConfig
)

type appConfig struct {
	Env  string `yaml:"env"`
	Name string `yaml:"name"`
	Log  struct {
		FilePath         string `mapstructure:"path"`
		FileMaxSize      int    `mapstructure:"maxsize"`
		BackUpFileMaxAge int    `mapstructure:"max_age"`
	}
	Pagination struct {
		DefaultPageSize int `mapstructure:"default_page_size"`
		MaxPageSize     int `mapstructure:"max_page_size"`
	}
}

type databaseConfig struct {
	Type        string        `mapstructure:"type"`
	DSN         string        `mapstructure:"dsn"`
	MaxOpenConn int           `mapstructure:"maxopen"`
	MaxIdleConn int           `mapstructure:"maxidle"`
	MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
}

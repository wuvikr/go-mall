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
}

type databaseConfig struct {
	Type        string        `mapstructure:"type"`
	DSN         string        `mapstructure:"dsn"`
	MaxOpenConn int           `mapstructure:"maxopen"`
	MaxIdleConn int           `mapstructure:"maxidle"`
	MaxLifeTime time.Duration `mapstructure:"maxlifetime"`
}

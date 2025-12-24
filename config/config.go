// config.go 定义配置对象
package config

import "time"

var (
	App      *appConfig
	Database *databaseConfig
	Redis    *redisConfig
)

type appConfig struct {
	Env  string `yaml:"env"`
	Name string `yaml:"name"`
	Log  struct {
		FilePath         string `mapstructure:"path"`
		FileMaxSize      int    `mapstructure:"max_size"`
		BackUpFileMaxAge int    `mapstructure:"max_age"`
	}
	Pagination struct {
		DefaultPageSize int `mapstructure:"default_page_size"`
		MaxPageSize     int `mapstructure:"max_page_size"`
	}
}

type databaseConfig struct {
	Type   string          `mapstructure:"type"`
	Master DBConnectOption `mapstructure:"master"`
	Slave  DBConnectOption `mapstructure:"slave"`
}

type DBConnectOption struct {
	DSN          string        `mapstructure:"dsn"`
	MaxOpenConns int           `mapstructure:"maxopen"`
	MaxIdleConns int           `mapstructure:"maxidle"`
	MaxLifeTime  time.Duration `mapstructure:"maxlifetime"`
}

type redisConfig struct {
	Addr     string `mapstructure:"host"`
	UserName string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

package config

import (
	"bytes"
	"embed"
	"os"
	"time"

	"github.com/spf13/viper"
)

//go:embed *.yaml
var configs embed.FS

func init() {
	env := os.Getenv("ENV")
	vp := viper.New()

	// 根据环境读取配置文件
	configFileStream, err := configs.ReadFile("application." + env + ".yaml")
	if err != nil {
		panic(err)
	}

	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewBuffer(configFileStream))
	if err != nil {
		panic(err)
	}

	vp.UnmarshalKey("app", &App)
	vp.UnmarshalKey("database", &Database)

	Database.MaxLifeTime = time.Second
}

// bootstrap.go 加载配置文件，把配置解析到配置对象
package config

import (
	"bytes"
	"embed"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

//go:embed *.yaml
var configs embed.FS

func init() {
	env, exist := os.LookupEnv("ENV")
	if !exist {
		panic("ENV 不存在，请检查环境变量！")
	}
	vp := viper.New()

	// 根据环境读取配置文件
	configFileStream, err := configs.ReadFile("application." + env + ".yaml")
	if err != nil {
		panic(fmt.Errorf("无法从 embed.FS 读取配置: %w", err))
	}

	vp.SetConfigType("yaml")
	err = vp.ReadConfig(bytes.NewReader(configFileStream))
	if err != nil {
		panic(err)
	}

	vp.UnmarshalKey("app", &App)
	vp.UnmarshalKey("database", &Database)

}

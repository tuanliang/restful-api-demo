package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

// 如何配置映射成Config对象

// 从Toml格式的配置文件加载配置
func LoadConfigFromToml(filePath string) error {
	config = NewDefaultConfig()
	// 读取toml格式的配置
	_, err := toml.DecodeFile(filePath, config)
	if err != nil {
		return fmt.Errorf("load config form file error %s", err)
	}

	return nil
}

// 从环境变量加载配置
func LoadConfigFromEnv() error {
	config = NewDefaultConfig()
	return env.Parse(config)
}

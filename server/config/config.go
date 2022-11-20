package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

var config *Config

type Config struct {
	// 服务器配置
	Server *ServerConfig `yaml:"server"`
	// 数据库配置
	Database *DatabaseConfig `yaml:"database"`
}

func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	var configFile string
	var err error

	config = new(Config)
	if configFile, err = GetConfigFile(); err != nil {
		return nil, err
	}

	// If the file exists,read the configuration from the file.
	if _, err := os.Stat(configFile); err == nil {
		cfb, err := os.ReadFile(configFile)
		if err != nil {
			return nil, err
		}
		if err = yaml.Unmarshal(cfb, config); err != nil {
			return nil, err
		}
	}

	if config.Server == nil {
		config.Server = new(ServerConfig)
	}
	if config.Database == nil {
		config.Database = new(DatabaseConfig)
	}

	if err := config.Server.getConfig(); err != nil {
		return nil, err
	}
	config.Database.setDefault()
	return config, nil
}

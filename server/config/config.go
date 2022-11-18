package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

var config *Config

type Config struct {
	// 服务器配置
	Server *ServerConfig `yaml:"server"`
}

func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	var configFile string
	var err error

	config = new(Config)
	server := new(ServerConfig)
	config.Server = server
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

	if err := server.getConfig(); err != nil {
		return nil, err
	}
	return config, nil
}

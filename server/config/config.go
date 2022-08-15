package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

var config *Config

type Config struct {
	Server *ServerConfig `yaml:"server"`
}

func GetConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	config = new(Config)
	server := new(ServerConfig)
	config.Server = server
	opt := GetOptions()

	// If the file exists,read the configuration from the file.
	if _, err := os.Stat(opt.ConfigFile); err == nil {
		cfb, err := os.ReadFile(opt.ConfigFile)
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

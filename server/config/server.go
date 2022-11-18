package config

import (
	"path"
	"path/filepath"
)

type ServerConfig struct {
	Listens []*ListenConfig `yaml:"listens"`
	// assets file directory
	Assets *string `yaml:"assets"`
	// HomePage
	HomePage *string `yaml:"homePage"`
}

func (config *ServerConfig) getConfig() (err error) {
	if err := config.correct(); err != nil {
		return err
	}
	if err := config.setDefault(); err != nil {
		return err
	}
	if err := config.checkConfig(); err != nil {
		return err
	}
	return nil
}

// correct: 对用户配置的参数进行规范化处理，比如单词改为全大写，文件的相对路径改为绝对路径等
func (config *ServerConfig) correct() (err error) {
	for _, item := range config.Listens {
		if err := item.correct(); err != nil {
			return err
		}
	}
	if config.Assets != nil {
		if *config.Assets, err = filepath.Abs(*config.Assets); err != nil {
			return err
		}
	}
	if config.HomePage != nil {
		if *config.HomePage, err = filepath.Abs(*config.HomePage); err != nil {
			return err
		}
	}
	return nil
}

// setDefault: 将用户没有配置的参数，设置为默认值
func (config *ServerConfig) setDefault() error {
	if len(config.Listens) == 0 {
		config.Listens = append(config.Listens, new(ListenConfig))
	}
	for _, item := range config.Listens {
		item.setDefault()
	}
	if config.Assets == nil {
		*config.Assets, _ = filepath.Abs("assets")
	}
	if config.HomePage == nil {
		*config.HomePage = path.Join(*config.Assets, "index.html")
	}
	return nil
}

func (config *ServerConfig) checkConfig() error {
	for _, item := range config.Listens {
		if err := item.checkConfig(); err != nil {
			return err
		}
	}
	return nil
}

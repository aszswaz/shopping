package config

import "path/filepath"

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	// 数据库地址
	Url *string `yaml:"url"`
}

func (database *DatabaseConfig) setDefault() {
	if database.Url == nil {
		database.Url = new(string)
		*database.Url, _ = filepath.Abs("shopping.db")
	}
}

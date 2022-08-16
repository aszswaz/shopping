package config

import (
	"errors"
	"giea.aszswaz.cn/aszswaz/shopping/utils"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ServerConfig struct {
	// SocketType: [TCP], [UNIX], [TCP, UNIX]
	SocketType []string `yaml:"socketType"`
	// Bind address
	Address string `yaml:"address"`
	// Unix domain socket path
	SocketPath string `yaml:"socketPath"`
	// Bind tcp port
	Port uint16 `yaml:"port"`
	// static file directory
	Static string `yaml:"static"`
	// home page
	Home string `yaml:"home"`
}

func (config *ServerConfig) getConfig() (err error) {
	opt := GetOptions()

	if opt.Address != "" {
		config.Address = opt.Address
	}
	if opt.Port != 0 {
		config.Port = opt.Port
	}
	if opt.Static != "" {
		config.Static = opt.Static
	}

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

// correct: Specification of user-configured parameters.
func (config *ServerConfig) correct() (err error) {
	if len(config.SocketType) > 0 {
		for index, item := range config.SocketType {
			config.SocketType[index] = strings.ToUpper(item)
		}
	}

	ss := utils.StringSlice{Slice: config.SocketType}
	config.SocketType = ss.Deduplication().Slice
	if ss.Contain("UNIX") && config.SocketPath != "" {
		if socketPath, err := filepath.Abs(config.SocketPath); err == nil {
			config.SocketPath = socketPath
		}
	}

	if config.Static != "" {
		if config.Static, err = filepath.Abs(config.Static); err != nil {
			return err
		}
	}
	if config.Home != "" {
		if config.Home, err = filepath.Abs(config.Home); err != nil {
			return err
		}
	}
	return nil
}

func (config *ServerConfig) setDefault() (err error) {
	if len(config.SocketType) == 0 {
		config.SocketType = append(config.SocketType, "TCP")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	ss := utils.StringSlice{Slice: config.SocketType}
	hasTcp := ss.Contain("TCP")
	hasUnix := ss.Contain("UNIX")

	if hasTcp && config.Address == "" {
		config.Address = "localhost"
	}
	if hasTcp && config.Port == 0 {
		config.Port = 8080
	}
	if hasUnix && config.SocketPath == "" {
		config.SocketPath = path.Join(cwd, "shopping.sock")
	}
	if config.Static == "" {
		config.Static = path.Join(cwd, "static")
	}
	if config.Home == "" {
		config.Home = path.Join(config.Static, "index.html")
	}
	return nil
}

func (config *ServerConfig) checkConfig() (err error) {
	if config.Address != "" && config.Address != "localhost" {
		if ip := net.ParseIP(config.Address); ip == nil {
			return errors.New("invalid bind address")
		}
	}
	if config.Port == 0 {
		return errors.New("invalid TCP port")
	}
	for _, item := range config.SocketType {
		if item != "TCP" && item != "UNIX" {
			return errors.New("unknown socket type:" + item)
		}
	}
	return err
}

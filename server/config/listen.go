package config

import (
	"errors"
	"net"
	"path/filepath"
)

type ListenConfig struct {
	// SocketType: socket 类型，TCP 或 unix 域套接字
	SocketType *string `yaml:"socketType"`
	// Address: TCP 的绑定地址或 unix 域套接字的路径
	Address *string `yaml:"address"`
	// Port: TCP 端口号
	Port *uint16 `yaml:"port"`
}

func (listen *ListenConfig) correct() (err error) {
	if *listen.SocketType == "UNIX" {
		if *listen.Address, err = filepath.Abs(*listen.Address); err != nil {
			return err
		}
	}
	return nil
}

func (listen *ListenConfig) setDefault() {
	if listen.SocketType == nil {
		listen.SocketType = new(string)
		*listen.SocketType = "TCP"
	}
	if *listen.SocketType == "TCP" {
		if listen.Address == nil {
			listen.Address = new(string)
			*listen.Address = "127.0.0.1"
		}
		if listen.Port == nil {
			listen.Port = new(uint16)
			*listen.Port = 8080
		}
	}
	if *listen.SocketType == "UNIX" && listen.Address == nil {
		listen.Address = new(string)
		*listen.Address, _ = filepath.Abs("shopping.sock")
	}
}

func (listen *ListenConfig) checkConfig() error {
	if *listen.SocketType == "TCP" {
		if ip := net.ParseIP(*listen.Address); ip == nil {
			return errors.New("invalid bind address")
		}
	} else if *listen.SocketType != "UNIX" {
		return errors.New("invalid socket type: " + *listen.SocketType)
	}
	return nil
}

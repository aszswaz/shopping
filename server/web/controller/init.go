package controller

import (
	"giea.aszswaz.cn/aszswaz/shopping/config"
	"giea.aszswaz.cn/aszswaz/shopping/hooks"
	"log"
)

var (
	serverCfg *config.ServerConfig
)

func init() {
	configObj, err := config.GetConfig()
	if err != nil {
		log.Println(err)
		hooks.ExitHandler(1)
	}
	serverCfg = configObj.Server
}

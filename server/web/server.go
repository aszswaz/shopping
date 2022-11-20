package web

import (
	"fmt"
	"giea.aszswaz.cn/aszswaz/shopping/config"
	"giea.aszswaz.cn/aszswaz/shopping/hooks"
	"giea.aszswaz.cn/aszswaz/shopping/web/controller"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	g         errgroup.Group
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

func Start() (err error) {
	router := routerManager()
	for _, item := range serverCfg.Listens {
		// 创建一个局部变量，用于协程的入口函数
		listen := item
		switch *item.SocketType {
		case "TCP":
			g.Go(func() error { return listenTCP(listen, router) })
		case "UNIX":
			g.Go(func() error { return listenUNIX(listen, router) })
		}
	}
	hooks.Register(exitHandler)
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// exitHandler: When the program exits, clean up the occupied resources.
func exitHandler() {
	cfg, _ := config.GetConfig()
	serverCfg := cfg.Server
	for _, item := range serverCfg.Listens {
		if *item.SocketType == "UNIX" {
			if _, err := os.Stat(*item.Address); err == nil {
				if err = os.Remove(*item.Address); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}

// listenTCP: listen http port
func listenTCP(listen *config.ListenConfig, router *gin.Engine) error {
	var host string
	if strings.Contains(*listen.Address, ":") {
		host = "[" + *listen.Address + "]" + ":" + strconv.FormatUint(uint64(*listen.Port), 10)
	} else {
		host = *listen.Address + ":" + strconv.FormatUint(uint64(*listen.Port), 10)
	}
	return router.Run(host)
}

// listenUNIX: listen unix domain socket
func listenUNIX(listen *config.ListenConfig, router *gin.Engine) error {
	if err := os.Remove(*listen.Address); err != nil && !os.IsNotExist(err) {
		return err
	}
	return router.RunUnix(*listen.Address)
}

// routerManager: 注册 URL 对应的处理函数
func routerManager() *gin.Engine {
	router := gin.Default()
	router.NoRoute(controller.Html404)
	router.GET("/", controller.Home)
	router.GET("/404.html", controller.Html404)
	router.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	router.Static("/assets", *serverCfg.Assets)

	user := router.Group("user")
	// 用户的登陆和注册接口
	user.POST("login", controller.Login)
	user.POST("register", controller.Register)
	// TODO: 获取用户信息的接口
	// TODO: 获取商店信息的接口
	return router
}

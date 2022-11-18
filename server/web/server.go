package web

import (
	"fmt"
	"giea.aszswaz.cn/aszswaz/shopping/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	g         errgroup.Group
	serverCfg *config.ServerConfig
)

func init() {
	configObj, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
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
	go signalExit()
	defer exitHandler()
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

// signalExit: Exit signal processing.
func signalExit() {
	sigCha := make(chan os.Signal)
	signal.Notify(sigCha, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	<-sigCha
	exitHandler()
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
	os.Exit(0)
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

// routerManager: Registration function management for HTTP requests.
func routerManager() *gin.Engine {
	router := gin.Default()
	router.NoRoute(html404)
	router.GET("/", home)
	router.GET("/404.html", html404)
	router.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	router.Static("/assets", *serverCfg.Assets)
	return router
}

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
	for _, sockType := range serverCfg.SocketType {
		switch sockType {
		case "TCP":
			g.Go(func() error { return listenTCP(router) })
		case "UNIX":
			g.Go(func() error { return listenUNIX(router) })
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
	if _, err := os.Stat(serverCfg.SocketPath); err == nil {
		if err = os.Remove(serverCfg.SocketPath); err != nil {
			//goland:noinspection GoUnhandledErrorResult
			fmt.Fprintln(os.Stderr, err)
		}
	}
	os.Exit(0)
}

// listenTCP: listen http port
func listenTCP(router *gin.Engine) error {
	var host string
	if strings.Contains(serverCfg.Address, ":") {
		host = "[" + serverCfg.Address + "]" + ":" + strconv.FormatUint(uint64(serverCfg.Port), 10)
	} else {
		host = serverCfg.Address + ":" + strconv.FormatUint(uint64(serverCfg.Port), 10)
	}
	return router.Run(host)
}

// listenUNIX: listen unix domain socket
func listenUNIX(router *gin.Engine) error {
	if err := os.Remove(serverCfg.SocketPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return router.RunUnix(serverCfg.SocketPath)
}

// routerManager: Registration function management for HTTP requests.
func routerManager() *gin.Engine {
	router := gin.Default()
	router.NoRoute(html404)
	router.GET("/", home)
	router.GET("/404.html", html404)
	router.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	router.Static("/static", serverCfg.Static)
	return router
}

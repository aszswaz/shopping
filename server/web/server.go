package web

import (
	"fmt"
	"giea.aszswaz.cn/aszswaz/shopping/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func Start() (err error) {
	configObj, err := config.GetConfig()
	if err != nil {
		return err
	}
	serverConfig := configObj.Server
	var exits [2]chan int8
	router := routerManager()
	for _, sockType := range serverConfig.SocketType {
		switch sockType {
		case "TCP":
			go listenTCP(serverConfig, router, exits[0])
		case "UNIX":
			go listenUNIX(serverConfig, router, exits[1])
		}
	}
	go signalExit()
	// Waiting for the exit signal from the service.
	for _, ch := range exits {
		<-ch
	}
	exitHandler()
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

// routerManager: Registration function management for HTTP requests.
func routerManager() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return router
}

// listenTCP: listen http port
func listenTCP(serverConfig *config.ServerConfig, router *gin.Engine, exit chan int8) {
	var host string
	if strings.Contains(serverConfig.Address, ":") {
		host = "[" + serverConfig.Address + "]" + ":" + strconv.FormatUint(uint64(serverConfig.Port), 10)
	} else {
		host = serverConfig.Address + ":" + strconv.FormatUint(uint64(serverConfig.Port), 10)
	}

	err := router.Run(host)
	if err != nil {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, err)
	}
	exit <- 1
}

// listenUNIX: listen unix domain socket
func listenUNIX(configObj *config.ServerConfig, router *gin.Engine, exit chan int8) {
	if err := os.Remove(configObj.SocketPath); err != nil && !os.IsNotExist(err) {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err := router.RunUnix(configObj.SocketPath)
	if err != nil {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, err)
	}
	exit <- 1
}

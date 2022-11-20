package hooks

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	hooks []func()
)

func init() {
	go signalExit()
}

// Register 注册程序退出事件的钩子
func Register(hook func()) {
	for _, item := range hooks {
		if &hook == &item {
			return
		}
	}
	hooks = append(hooks, hook)
}

// signalExit 收到退出信号
func signalExit() {
	sigCha := make(chan os.Signal)
	signal.Notify(sigCha, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	<-sigCha
	ExitHandler(1)
}

// ExitHandler 执行程序退出钩子，并退出程序
func ExitHandler(code int) {
	for _, item := range hooks {
		item()
	}
	os.Exit(code)
}

package hooks

import (
	"os"
	"os/signal"
	"reflect"
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
	// 获取函数的地址，根据函数的地址进行去重
	funcPtr1 := reflect.ValueOf(hook).Pointer()
	for i := 0; i < len(hooks); i++ {
		funcPtr2 := reflect.ValueOf(hooks[i]).Pointer()
		if funcPtr1 == funcPtr2 {
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

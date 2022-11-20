package main

import (
	"giea.aszswaz.cn/aszswaz/shopping/hooks"
	"giea.aszswaz.cn/aszswaz/shopping/web"
	"log"
)

func main() {
	code := 0
	if err := web.Start(); err != nil {
		log.Println(err)
		code = 1
	}
	hooks.ExitHandler(code)
}

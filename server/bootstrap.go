package main

import (
	"fmt"
	"giea.aszswaz.cn/aszswaz/shopping/web"
	"os"
)

func main() {
	if err := web.Start(); err != nil {
		//goland:noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

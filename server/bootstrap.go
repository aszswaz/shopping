package main

import (
	"giea.aszswaz.cn/aszswaz/shopping/web"
	"log"
)

func main() {
	if err := web.Start(); err != nil {
		log.Fatalln(err)
	}
}

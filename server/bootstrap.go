package main

import (
	"fmt"
	"giea.aszswaz.cn/aszswaz/shopping/web"
	"log"
)

func main() {
	fmt.Println("Hello World")
	if err := web.Start(); err != nil {
		log.Fatalln(err)
	}
}

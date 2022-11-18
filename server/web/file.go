package web

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func home(context *gin.Context) {
	if _, err := os.Stat(*serverCfg.HomePage); err == nil {
		context.File(*serverCfg.HomePage)
	} else {
		html404(context)
	}
}

func html404(context *gin.Context) {
	pagePath := path.Join(*serverCfg.Assets, "404.html")
	page, err := os.ReadFile(pagePath)
	if err != nil {
		page = []byte("<!DOCTYPE html>" +
			"<html lang=\"en\">" +
			"<head>" +
			"<meta charset=\"UTF-8\">" +
			"<title>Not Found</title>" +
			"</head>" +
			"<body>" +
			"<h1 style=\"text-align: center\">404 Not Found</h1>" +
			"</body>" +
			"</html>")
	}
	context.Status(http.StatusNotFound)
	context.Header("Content-Type", gin.MIMEHTML)
	context.Header("Content-Length", strconv.Itoa(len(page)))
	if _, err := context.Writer.Write(page); err != nil {
		log.Println(err)
	}
}

package web

import (
	"github.com/gin-gonic/gin"
)

func home(context *gin.Context) {
	context.File(serverCfg.Home)
}

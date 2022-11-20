package web

import "github.com/gin-gonic/gin"

// LoginInterceptor 对请求的登陆状态进行校验
func LoginInterceptor(context *gin.Context) {
	// TODO: 校验请求中的 COOKIE
}

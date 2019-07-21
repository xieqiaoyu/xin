package http

import (
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin"
)

var httpEngine *gin.Engine

//RegistRouterAndMiddlewareFun 给gin 引擎注册中间键的函数
type RegistRouterAndMiddlewareFunc func(r *gin.Engine)

//Engine 获取http server 对象
func Engine(RegistRouterAndMiddleware RegistRouterAndMiddlewareFunc) *gin.Engine {
	var mode string
	switch xin.Mode() {
	case xin.Dev:
		mode = "debug"
	case xin.Test:
		mode = "test"
	case xin.Release:
		mode = "release"
	}
	gin.SetMode(mode)
	httpEngine = gin.New()
	RegistRouterAndMiddleware(httpEngine)
	return httpEngine
}

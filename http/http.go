package http

import (
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin"
)

//RegistRouterAndMiddlewareFun  handle to register router and middleware for gin engine
type RegisterRouterAndMiddlewareFunc func(r *gin.Engine)

//Engine wrapper to get gin engine instanee
func Engine() *gin.Engine {
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
	httpEngine := gin.New()
	return httpEngine
}

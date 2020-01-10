package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

//ListenAndServe gradeful start http server
func ListenAndServe(r *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	xlog.WriteInfo("Http server working on %s", addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			xlog.WriteError("Ooops! %s", err)
			//TODO: return error instead of exit
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//TODO: return when go routine exit
	xlog.WriteInfo("Shutdown Server ...")
	ctx := context.Background()
	if err := srv.Shutdown(ctx); err != nil {
		xlog.WriteError("Server Shutdown: %s", err)
		os.Exit(1)
	}
	xlog.WriteInfo("Server exited")
}

package cmd

import (
	"github.com/spf13/cobra"

	"context"
	"github.com/xieqiaoyu/xin"
	httpserver "github.com/xieqiaoyu/xin/http"
	xlog "github.com/xieqiaoyu/xin/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var registRouterFunc httpserver.RegistRouterAndMiddlewareFunc

//UseRouterRegister 指定使用的路由注册函数
func UseRouterRegister(register httpserver.RegistRouterAndMiddlewareFunc) {
	registRouterFunc = register
}

//ServerCmd 启动服务器的命令
func HttpServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "http server",
		Long:  `control http server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			if ConfigFileToUse != "" {
				xin.SetConfigFile(ConfigFileToUse)
			}
			if err := xin.LoadConfig(); err != nil {
				xlog.WriteError("%s", err)
				os.Exit(1)
			}
			r := httpserver.Engine(registRouterFunc)
			srv := &http.Server{
				Addr:    xin.Config().GetString("http.listen"),
				Handler: r,
			}
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					xlog.WriteError("Ooops! %s", err)
					os.Exit(1)
				}
			}()
			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			xlog.WriteInfo("Shutdown Server ...")
			ctx := context.Background()
			if err := srv.Shutdown(ctx); err != nil {
				xlog.WriteError("Server Shutdown: %s", err)
				os.Exit(1)
			}
			xlog.WriteInfo("Server exiting")
		},
	}
}

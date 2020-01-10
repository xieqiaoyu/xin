package cmd

import (
	"github.com/spf13/cobra"

	"github.com/xieqiaoyu/xin"
	xcmd "github.com/xieqiaoyu/xin/cmd"
	xhttp "github.com/xieqiaoyu/xin/http"
	xlog "github.com/xieqiaoyu/xin/log"
	"os"
)

var registerRouterAndMiddare xhttp.RegisterRouterAndMiddlewareFunc

//UseRouterRegister 指定使用的路由注册函数
func UseRouterRegister(handle xhttp.RegisterRouterAndMiddlewareFunc) {
	registerRouterAndMiddare = handle
}

//Cmd 启动服务器的命令
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "http server",
		Long:  `control http server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := xcmd.ConfigInit(); err != nil {
				xlog.WriteError("%s", err)
				os.Exit(1)
			}
			addr := xin.Config().GetString("http.listen")
			if addr == "" {
				addr = ":8080"
			}
			r := xhttp.Engine()
			registerRouterAndMiddare(r)
			xhttp.ListenAndServe(r, addr)
		},
	}
}

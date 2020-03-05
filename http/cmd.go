package http

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	"net/http"
	"os"
)

//ServerInterface a server can provide http server
type ServerInterface interface {
	// provide the http server service
	GetHTTPServer() (*http.Server, error)
}

//InitializeServerFunc function init http Server  gives the posibility for dependence inject
type InitializeServerFunc func() (ServerInterface, error)

//NewHTTPCmd Get a cobra command start http server
func NewHTTPCmd(getServer InitializeServerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "http server",
		Long:  `control http server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			server, err := getServer()
			if err != nil {
				xlog.Errorf("Init server fail %s", err)
				os.Exit(1)
			}
			httpServer, err := server.GetHTTPServer()
			if err != nil {
				xlog.Errorf("Fail to get http Server : %s", err)
				return
			}

			addr := httpServer.Addr
			if addr == "" {
				addr = ":http"
			}

			xlog.Infof("Http server working on %s", addr)
			go func() {
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					xlog.Errorf("Ooops! %s", err)
					os.Exit(1)
				}
			}()

			xin.WaitForQuitSignal()

			xlog.Infof("Shutdown Server ...")
			ctx := context.Background()
			if err := httpServer.Shutdown(ctx); err != nil {
				xlog.Errorf("Server Shutdown: %s", err)
				os.Exit(1)
			}
			xlog.Infof("Server exited")
		},
	}
}

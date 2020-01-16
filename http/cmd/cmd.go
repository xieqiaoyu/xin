package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	GetHttpServer() *http.Server
}

type InitializeServerFunc func(config *xin.Config) (Server, error)

//NewHttpCmd Get a cobra command start http server
func NewHttpCmd(config *xin.Config, initServer InitializeServerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "http server",
		Long:  `control http server behavior`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := config.Init(); err != nil {
				xlog.WriteError("%s", err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			server, err := initServer(config)
			if err != nil {
				xlog.WriteError("Init server fail %s", err)
				os.Exit(1)
			}
			httpServer := server.GetHttpServer()

			addr := httpServer.Addr
			if addr == "" {
				addr = ":http"
			}

			xlog.WriteInfo("Http server working on %s", addr)
			go func() {
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					xlog.WriteError("Ooops! %s", err)
					os.Exit(1)
				}
			}()

			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			xlog.WriteInfo("Shutdown Server ...")
			ctx := context.Background()
			if err := httpServer.Shutdown(ctx); err != nil {
				xlog.WriteError("Server Shutdown: %s", err)
				os.Exit(1)
			}
			xlog.WriteInfo("Server exited")
		},
	}
}

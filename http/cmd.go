package http

import (
	"context"
	"github.com/spf13/cobra"
	xlog "github.com/xieqiaoyu/xin/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ServerInterface interface {
	GetHttpServer() *http.Server
}

type InitializeServerFunc func() (ServerInterface, error)

//NewHttpCmd Get a cobra command start http server
func NewHttpCmd(getServer InitializeServerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "http server",
		Long:  `control http server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			server, err := getServer()
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

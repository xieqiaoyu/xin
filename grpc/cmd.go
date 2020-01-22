package grpc

import (
	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

type ServerInterface interface {
	GetGrpcServer() (*grpc.Server, error)
	GetNetListener() (net.Listener, error)
}

type InitializeServerFunc func() (ServerInterface, error)

func NewGrpcCmd(getServer InitializeServerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "grpc",
		Short: "start grpc service",
		Long:  `control grpc server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			server, err := getServer()
			if err != nil {
				xlog.WriteError("Init server fail %s", err)
				os.Exit(1)
			}

			s, err := server.GetGrpcServer()
			if err != nil {
				xlog.WriteError("Get Gprc server fail %s", err)
				os.Exit(1)
			}
			lis, err := server.GetNetListener()
			if err != nil {
				xlog.WriteError("Get Gprc net listener fail %s", err)
				os.Exit(1)
			}
			addr := lis.Addr()

			go func() {
				xlog.WriteInfo("Grpc server working on %s/%s", addr.Network(), addr.String())
				err := s.Serve(lis)
				if err != nil {
					xlog.WriteError("%s", err)
					os.Exit(1)
				}
			}()

			xin.WaitForQuitSignal()

			xlog.WriteInfo("Shutdown Server ...")
			s.GracefulStop()
			xlog.WriteInfo("Server exited")
		},
	}
}

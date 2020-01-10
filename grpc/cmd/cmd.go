package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xcmd "github.com/xieqiaoyu/xin/cmd"
	xgrpc "github.com/xieqiaoyu/xin/grpc"
	xlog "github.com/xieqiaoyu/xin/log"
	"google.golang.org/grpc"
)

var grpcServerRegister xgrpc.RegistServerFunc
var grpcServerOpts []grpc.ServerOption

//UseGRPCServerRegister UseGRPCServerRegister
func UseServerRegister(register xgrpc.RegistServerFunc) {
	grpcServerRegister = register
}

//SetGRPCServerOpts SetGRPCServerOpts
func SetServerOpts(opts []grpc.ServerOption) {
	grpcServerOpts = opts
}

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "grpc",
		Short: "start grpc service",
		Long:  `control grpc server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := xcmd.ConfigInit(); err != nil {
				xlog.WriteError("%s", err)
				os.Exit(1)
			}

			s := xgrpc.Server(grpcServerOpts...)

			if grpcServerRegister != nil {
				grpcServerRegister(s)
			}
			addr := xin.Config().GetString("grpc.listen")
			if addr == "" {
				addr = ":50051"
			}

			go func() {
				xlog.WriteInfo("Grpc server working on %s", addr)
				err := xgrpc.ServeTCP(s, addr)
				if err != nil {
					xlog.WriteError("%s", err)
					os.Exit(1)
				}
			}()

			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			xlog.WriteInfo("Shutdown Server ...")
			s.GracefulStop()
			xlog.WriteInfo("Server exited")
		},
	}
}

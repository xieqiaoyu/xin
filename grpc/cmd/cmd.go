package cmd

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xcmd "github.com/xieqiaoyu/xin/cmd"
	xgrpc "github.com/xieqiaoyu/xin/grpc"
	xlog "github.com/xieqiaoyu/xin/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var grpcServerRegister xgrpc.RegistServerFunc
var grpcServerOpts []grpc.ServerOption

//UseGRPCServerRegister UseGRPCServerRegister
func UseGRPCServerRegister(register xgrpc.RegistServerFunc) {
	grpcServerRegister = register
}

//SetGRPCServerOpts SetGRPCServerOpts
func SetGRPCServerOpts(opts []grpc.ServerOption) {
	grpcServerOpts = opts
}

func GrpcServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "grpc",
		Short: "start grpc service",
		Long:  `control grpc server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := xcmd.ConfigInit(); err != nil {
				xlog.WriteError("%s", err)
				os.Exit(1)
			}
			s := grpc.NewServer(grpcServerOpts...)
			if grpcServerRegister != nil {
				grpcServerRegister(s)
			}
			// if in development we enable grpc reflection
			if xin.Mode() == xin.Dev {
				reflection.Register(s)
			}
			addr := xgrpc.Addr()
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				xlog.WriteError("failed to listen: %v", err)
				os.Exit(1)
			}

			xlog.WriteInfo("Grpc server working on %s", addr)
			go func() {
				if err := s.Serve(lis); err != nil {
					xlog.WriteError("failed to serve: %v", err)
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

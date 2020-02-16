package grpc

import (
	"github.com/spf13/cobra"
	"github.com/xieqiaoyu/xin"
	xlog "github.com/xieqiaoyu/xin/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

//ServerInterface Grpc service interface
type ServerInterface interface {
	//GetGrpcServer get the grpc server
	GetGrpcServer() (*grpc.Server, error)
	//GetNetListener get grpc Net Listener
	GetNetListener() (net.Listener, error)
}

//InitializeServerFunc function to Initialize grpc server
type InitializeServerFunc func() (ServerInterface, error)

//NewGrpcCmd get a cobra command to run grpc service
func NewGrpcCmd(getServer InitializeServerFunc) *cobra.Command {
	return &cobra.Command{
		Use:   "grpc",
		Short: "start grpc service",
		Long:  `control grpc server behavior`,
		Run: func(cmd *cobra.Command, args []string) {
			server, err := getServer()
			if err != nil {
				xlog.Errorf("Init server fail %s", err)
				os.Exit(1)
			}

			s, err := server.GetGrpcServer()
			if err != nil {
				xlog.Errorf("Get Gprc server fail %s", err)
				os.Exit(1)
			}
			lis, err := server.GetNetListener()
			if err != nil {
				xlog.Errorf("Get Gprc net listener fail %s", err)
				os.Exit(1)
			}
			addr := lis.Addr()

			go func() {
				xlog.Infof("Grpc server working on %s/%s", addr.Network(), addr.String())
				err := s.Serve(lis)
				if err != nil {
					xlog.Errorf("%s", err)
					os.Exit(1)
				}
			}()

			xin.WaitForQuitSignal()

			xlog.Infof("Shutdown Server ...")
			s.GracefulStop()
			xlog.Infof("Server exited")
		},
	}
}

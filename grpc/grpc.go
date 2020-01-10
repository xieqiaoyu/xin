package grpc

import (
	"github.com/xieqiaoyu/xin"
	"google.golang.org/grpc"
)

//RegistGRPCServerFunc RegistGRPCServerFunc
type RegistServerFunc func(*grpc.Server)

func Addr() string {
	addr := xin.Config().GetString("grpc.listen")
	if addr == "" {
		addr = ":50051"
	}
	return addr
}

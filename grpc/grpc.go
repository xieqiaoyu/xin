package grpc

import (
	"github.com/xieqiaoyu/xin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

//RegistServerFunc RegistServerFunc
type RegistServerFunc func(*grpc.Server)

func Server(opts ...grpc.ServerOption) *grpc.Server {
	s := grpc.NewServer(opts...)
	// if in development we enable grpc reflection
	if xin.Mode() == xin.Dev {
		reflection.Register(s)
	}
	return s
}

func ServeTCP(s *grpc.Server, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return xin.NewWrapEf("failed to listen: %w", err)
	}

	if err := s.Serve(lis); err != nil {
		return xin.NewWrapEf("failed to serve: %w", err)
	}
	return nil
}

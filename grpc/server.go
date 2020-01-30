package grpc

import (
	"github.com/xieqiaoyu/xin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

//RegistServerFunc RegistServerFunc
type RegistServerFunc func(*grpc.Server)

//ServerConfig config provide grpc server setting
type ServerConfig interface {
	//Get Grpc server listen setting network :tcp,udp  address
	GrpcListen() (network, address string)
}

//Server an implemention of ServerInterface
type Server struct {
	config              ServerConfig
	opts                []grpc.ServerOption
	env                 xin.Envirment
	registServerHandler RegistServerFunc
}

//NewServer Get a new grpc service server
func NewServer(config ServerConfig, env xin.Envirment, opts []grpc.ServerOption, registServerHandler RegistServerFunc) *Server {
	return &Server{
		config:              config,
		opts:                opts,
		env:                 env,
		registServerHandler: registServerHandler,
	}
}

//GetGrpcServer ServerInterface implement
func (s *Server) GetGrpcServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer(s.opts...)
	if s.env.Mode() == xin.DevMode {
		reflection.Register(grpcServer)
	}
	if s.registServerHandler != nil {
		s.registServerHandler(grpcServer)
	}
	return grpcServer, nil
}

//GetNetListener ServerInterface implement
func (s *Server) GetNetListener() (net.Listener, error) {
	network, addr := s.config.GrpcListen()
	if network == "" {
		network = "tcp"
	}
	return net.Listen(network, addr)
}

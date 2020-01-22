package grpc

import (
	"context"
	"fmt"
	"github.com/xieqiaoyu/xin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"strconv"
	"time"
)

//CustomError grpc CustomError interface
type CustomError interface {
	error
	GetCode() int
	GetMsg() string
}

//RegistServerFunc RegistServerFunc
type RegistServerFunc func(*grpc.Server)

type ServerConfig interface {
	GrpcListen() (network, address string)
}

type Server struct {
	config              ServerConfig
	opts                []grpc.ServerOption
	env                 xin.Envirment
	registServerHandler RegistServerFunc
}

func NewServer(config ServerConfig, env xin.Envirment, opts []grpc.ServerOption, registServerHandler RegistServerFunc) *Server {
	return &Server{
		config:              config,
		opts:                opts,
		env:                 env,
		registServerHandler: registServerHandler,
	}
}

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

func (s *Server) GetNetListener() (net.Listener, error) {
	network, addr := s.config.GrpcListen()
	if network == "" {
		network = "tcp"
	}
	return net.Listen(network, addr)
}

type UnaryContext struct {
	interceptors []UnaryChainServerInterceptor
	curIndex     int
	maxIndex     int
	grpcCtx      context.Context       // same as grpc.UnaryServerInterceptor ctx
	Req          interface{}           // same as grpc.UnaryServerInterceptor req
	Info         *grpc.UnaryServerInfo // same as grpc.UnaryServerInterceptor info
	Resp         interface{}           // same as UnaryServerInterceptor return value resp
	RespErr      error                 // same as UnaryServerInterceptor return value err
}

func (c *UnaryContext) Deadline() (deadline time.Time, ok bool) {
	return c.grpcCtx.Deadline()
}
func (c *UnaryContext) Err() error {
	return c.grpcCtx.Err()
}

func (c *UnaryContext) Done() <-chan struct{} {
	return c.grpcCtx.Done()
}

func (c *UnaryContext) Value(key interface{}) interface{} {
	return c.grpcCtx.Value(key)
}

func (c *UnaryContext) Next() {
	if c.curIndex < c.maxIndex {
		handler := c.interceptors[c.curIndex]
		c.curIndex++
		handler(c)
	}
}

func (c *UnaryContext) Abort() {
	c.curIndex = c.maxIndex
}

func (c *UnaryContext) IsAborted() bool {
	return c.curIndex >= c.maxIndex
}

type UnaryChainServerInterceptor func(c *UnaryContext)

//UnaryChainInterceptor  Option set the ServerChainInterceptor The first interceptor will be the outer most, while the last interceptor will be the inner most wrapper around the real call.
func UnaryChainInterceptor(interceptors ...UnaryChainServerInterceptor) grpc.ServerOption {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		i := interceptors
		i = append(i, func(uc *UnaryContext) {
			uc.Resp, uc.RespErr = handler(ctx, req)
		})

		c := &UnaryContext{
			interceptors: i,
			curIndex:     0,
			maxIndex:     len(i),
			grpcCtx:      ctx,
			Req:          req,
			Info:         info,
		}

		for !c.IsAborted() {
			c.Next()
		}
		return c.Resp, c.RespErr
	}
	return grpc.UnaryInterceptor(interceptor)
}

//UnaryCustomErrorRender if api return a custom Error render it into grpc trailer, this is an UnaryChainServerInterceptor
func UnaryCustomErrorRender(c *UnaryContext) {
	c.Next()
	if cusErr, ok := c.RespErr.(CustomError); ok {
		code := cusErr.GetCode()
		codeStr := strconv.Itoa(code)
		msg := cusErr.GetMsg()
		grpc.SetTrailer(c, metadata.Pairs("code", codeStr))
		grpc.SetTrailer(c, metadata.Pairs("info", msg))
		c.RespErr = status.Error(codes.Aborted, fmt.Sprintf("code %d, info:%s", code, msg))
		return
	}
}

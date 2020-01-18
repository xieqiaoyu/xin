package server

import (
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin"
	"net/http"
)

type Service interface {
	RegisterRouter(*gin.Engine)
}

type Server struct {
	//config  *xin.Config
	config  ServerConfig
	env     xin.Envirment
	service Service
}

type ServerConfig interface {
	HttpAddr() string
}

func (s *Server) GetHttpServer() *http.Server {
	//addr := s.config.Viper().GetString("http.listen")
	addr := s.config.HttpAddr()
	if addr == "" {
		addr = ":8080"
	}
	var mode string
	switch s.env.Mode() {
	case xin.DevMode:
		mode = "debug"
	case xin.TestMode:
		mode = "test"
	case xin.ReleaseMode:
		mode = "release"
	}
	gin.SetMode(mode)
	r := gin.New()
	if s.service != nil {
		s.service.RegisterRouter(r)
	}

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func NewServer(env xin.Envirment, config ServerConfig, service Service) *Server {
	return &Server{
		config:  config,
		env:     env,
		service: service,
	}
}

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/xieqiaoyu/xin"
	"net/http"
)

//Service http service interface
type Service interface {
	// register route and middleware into gin engine
	RegisterRouter(*gin.Engine)
}

//Server  Http server implement ServerInterface
type Server struct {
	config  ServerConfig
	env     xin.Envirment
	service Service
}

//ServerConfig config provide HTTP server setting
type ServerConfig interface {
	HTTPListen() string
}

//GetHTTPServer ServerInterface implement
func (s *Server) GetHTTPServer() *http.Server {
	addr := s.config.HTTPListen()
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

//NewServer Create a new HTTP server
func NewServer(env xin.Envirment, config ServerConfig, service Service) *Server {
	return &Server{
		config:  config,
		env:     env,
		service: service,
	}
}

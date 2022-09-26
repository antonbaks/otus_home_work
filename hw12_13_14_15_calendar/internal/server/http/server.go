package internalhttp

import (
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	logger  Logger
	cfg     Config
	server  *http.Server
	servMux *runtime.ServeMux
}

type Logger interface {
	Error(msg string)
}

type Config interface {
	GetHTTPHost() string
	GetHTTPGatewayPort() string
}

func NewServer(logger Logger, cfg Config, servMux *runtime.ServeMux) *Server {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		servMux: servMux,
	}
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.cfg.GetHTTPHost(), s.cfg.GetHTTPGatewayPort()),
		Handler: LoggingMiddleware(s.servMux),
	}

	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	return s.server.Close()
}

package internalhttp

import (
	"context"
	"net"
	"net/http"
)

type Server struct {
	logger Logger
	app    Application
	config Config
	Server *http.Server
}

type Logger interface { // TODO
}

type Application interface {
	PrintHello(w http.ResponseWriter, r *http.Request)
}

type Config interface {
	GetHTTPHost() string
	GetHTTPPort() string
}

func NewServer(logger Logger, app Application, config Config) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(s.app.PrintHello)
	mux.Handle("/hello", loggingMiddleware(finalHandler))

	server := &http.Server{
		Addr:    net.JoinHostPort(s.config.GetHTTPHost(), s.config.GetHTTPPort()),
		Handler: mux,
	}
	defer server.Close()
	s.Server = server

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx.Done()

	s.Server.Close()

	return nil
}

// TODO

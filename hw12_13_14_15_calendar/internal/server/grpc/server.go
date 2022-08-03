package internalgrpc

import (
	"context"
	"net"

	"github.com/antonbaks/otus_home_work/hw12_13_14_15_calendar/api/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Service Service
	ServMux *runtime.ServeMux
	logg    Logger
	cfg     Config
	serv    *grpc.Server
}

type Logger interface {
	Error(msg string)
}

type Config interface {
	GetGRPCPort() string
	GetHTTPHost() string
	GetHTTPType() string
	GetHTTPGatewayPort() string
}

func NewServer(service Service, logg Logger, cfg Config) *Server {
	return &Server{Service: service, logg: logg, cfg: cfg, ServMux: runtime.NewServeMux()}
}

func (s *Server) Start(ctx context.Context) error {
	address := net.JoinHostPort(s.cfg.GetHTTPHost(), s.cfg.GetGRPCPort())
	lsn, err := net.Listen(s.cfg.GetHTTPType(), address)
	if err != nil {
		return err
	}

	s.serv = grpc.NewServer(withServerUnaryInterceptor())
	pb.RegisterEventsServer(s.serv, &s.Service)

	go func() {
		s.serv.Serve(lsn)
	}()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	if err = pb.RegisterEventsHandler(ctx, s.ServMux, conn); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	s.serv.Stop()
}

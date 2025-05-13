package transport

import (
	"auth/internal/config"
	"auth/internal/service"
	transport "auth/internal/transport/grpc"
	"fmt"
	"strconv"

	"auth/pkg/logger"
	"context"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	cfg    config.Config
}

func NewServer(cfg config.Config, authService service.AuthService) *Server {
	srv := grpc.NewServer()

	transport.Register(srv, authService)

	return &Server{
		server: srv,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.cfg.Grpc.GRPCPort))
	if err != nil {
		return err
	}

	ctx := context.Background()
	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Starting gRPC server on port :%d", s.cfg.Grpc.GRPCPort))
	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}

package server

import (
	"fmt"
	"net"

	"github.com/koer/koer-module/pkg/config"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	cfg    config.GRPCConfig
}

func NewGRPCServer(cfg config.GRPCConfig, opts ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		server: grpc.NewServer(opts...),
		cfg:    cfg,
	}
}

func (s *GRPCServer) Server() *grpc.Server {
	return s.server
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("listening on port %d: %w", s.cfg.Port, err)
	}
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("serving grpc: %w", err)
	}
	return nil
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

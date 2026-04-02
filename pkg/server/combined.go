package server

import (
	"fmt"

	"github.com/koer/koer-module/pkg/config"
	"google.golang.org/grpc"
)

type CombinedServer struct {
	http *HTTPServer
	grpc *GRPCServer
}

func NewCombinedServer(cfg config.CombinedConfig, grpcOpts ...grpc.ServerOption) *CombinedServer {
	return &CombinedServer{
		http: NewHTTPServer(cfg.HTTP),
		grpc: NewGRPCServer(cfg.GRPC, grpcOpts...),
	}
}

func (s *CombinedServer) HTTP() *HTTPServer {
	return s.http
}

func (s *CombinedServer) GRPC() *GRPCServer {
	return s.grpc
}

func (s *CombinedServer) Start() error {
	errCh := make(chan error, 2)
	go func() { errCh <- s.http.Start() }()
	go func() { errCh <- s.grpc.Start() }()
	if err := <-errCh; err != nil {
		return fmt.Errorf("combined server error: %w", err)
	}
	return nil
}

func (s *CombinedServer) Stop() {
	s.grpc.Stop()
	_ = s.http.Shutdown()
}

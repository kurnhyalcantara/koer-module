package connection

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/koer/koer-module/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClientManager manages outgoing gRPC client connections to downstream
// services. It centralises dial-option setup (TLS / insecure), connection
// establishment, and lifecycle management so that every consumer service
// does not have to duplicate this orchestration.
//
// Service addresses are resolved in order:
//  1. Environment variable override: {UPPER_SNAKE_NAME}_ADDR (e.g. AUTH_SERVICE_ADDR)
//  2. Default from ServiceRegistry: localhost:{standardized_port}
//
// Usage:
//
//	cfg.GRPCClient.Services = []config.GRPCServiceTarget{
//		{Name: "auth-service", Enabled: true},
//		{Name: "user-service", Enabled: true},
//	}
//	mgr, err := connection.NewGRPCClientManager(cfg.GRPCClient)
//	defer mgr.Close()
//
//	authConn, err := mgr.Conn("auth-service")
//	authClient := authv1.NewAuthServiceClient(authConn)
type GRPCClientManager struct {
	conns map[string]*grpc.ClientConn
}

// NewGRPCClientManager creates a new manager and dials all enabled service
// targets declared in cfg.Services. If cfg.CertFile is empty, the manager
// uses insecure credentials; otherwise it loads TLS from the file.
// Targets with Enabled=false are silently skipped.
func NewGRPCClientManager(cfg config.GRPCClientConfig) (*GRPCClientManager, error) {
	opts, err := buildGRPCDialOpts(cfg.CertFile)
	if err != nil {
		return nil, err
	}

	mgr := &GRPCClientManager{conns: make(map[string]*grpc.ClientConn, len(cfg.Services))}

	for _, svc := range cfg.Services {
		if !svc.Enabled {
			log.Printf("[ grpc-client-manager ] skipping %s (disabled)", svc.Name)
			continue
		}

		addr := resolveServiceAddr(svc.Name)
		if addr == "" {
			mgr.Close()
			return nil, fmt.Errorf("grpc client: %s is not in ServiceRegistry and no env override set", svc.Name)
		}

		conn, err := grpc.NewClient(addr, opts...)
		if err != nil {
			mgr.Close()
			return nil, fmt.Errorf("grpc client: dial %s at %s: %w", svc.Name, addr, err)
		}

		mgr.conns[svc.Name] = conn
		log.Printf("[ grpc-client-manager ] connected to %s at %s", svc.Name, addr)
	}

	return mgr, nil
}

// Conn returns the established connection for the given service name.
// Returns an error if the service was not registered or was disabled.
func (m *GRPCClientManager) Conn(serviceName string) (*grpc.ClientConn, error) {
	conn, ok := m.conns[serviceName]
	if !ok {
		return nil, fmt.Errorf("grpc client: no connection for %q (not registered or disabled)", serviceName)
	}
	return conn, nil
}

// Close gracefully closes all tracked gRPC connections.
func (m *GRPCClientManager) Close() {
	for _, conn := range m.conns {
		conn.Close()
	}
}

func buildGRPCDialOpts(certFile string) ([]grpc.DialOption, error) {
	if certFile == "" {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return nil, fmt.Errorf("grpc client: load TLS credentials: %w", err)
	}
	return []grpc.DialOption{grpc.WithTransportCredentials(creds)}, nil
}

// resolveServiceAddr returns the address for a service by checking:
//  1. Env var override: {UPPER_SNAKE_NAME}_ADDR (e.g. AUTH_SERVICE_ADDR)
//  2. Default from ServiceRegistry: localhost:{standardized_port}
func resolveServiceAddr(name string) string {
	envKey := strings.ToUpper(strings.ReplaceAll(name, "-", "_")) + "_ADDR"
	if addr := os.Getenv(envKey); addr != "" {
		return addr
	}
	return config.DefaultServiceAddr(name)
}

package config

import "fmt"

// ServiceRegistry maps standardized service names to their default gRPC ports.
// Every service in the ecosystem gets a fixed port so all teams share the same
// local-development topology without per-service configuration.
//
// To add a new service, register it here once — every consumer picks it up
// automatically.
var ServiceRegistry = map[string]int{
	"auth-service":    9091,
	"user-service":    9092,
	"product-service": 9093,
}

// DefaultServiceAddr returns "localhost:{port}" for a registered service.
// Returns an empty string if the service is not in the registry.
func DefaultServiceAddr(name string) string {
	port, ok := ServiceRegistry[name]
	if !ok {
		return ""
	}
	return fmt.Sprintf("localhost:%d", port)
}

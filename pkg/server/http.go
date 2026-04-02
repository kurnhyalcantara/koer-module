package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/koer/koer-module/pkg/config"
)

type HTTPServer struct {
	app *fiber.App
	cfg config.HTTPConfig
}

func NewHTTPServer(cfg config.HTTPConfig) *HTTPServer {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
	return &HTTPServer{app: app, cfg: cfg}
}

func (s *HTTPServer) App() *fiber.App {
	return s.app
}

func (s *HTTPServer) Start() error {
	addr := fmt.Sprintf(":%d", s.cfg.Port)
	if err := s.app.Listen(addr); err != nil {
		return fmt.Errorf("starting http server: %w", err)
	}
	return nil
}

func (s *HTTPServer) Shutdown() error {
	if err := s.app.Shutdown(); err != nil {
		return fmt.Errorf("shutting down http server: %w", err)
	}
	return nil
}

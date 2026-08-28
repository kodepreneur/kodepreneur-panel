package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kodepreneur/agent/config"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(cfg *config.Config) *Server {
	router := NewRouter(cfg)

	httpServer := &http.Server{
		Handler:           router.Handler(),
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      600 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	return &Server{
		cfg:        cfg,
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	if s.cfg.Server.UseSocket {
		_ = os.Remove(s.cfg.Server.SocketPath)
		_ = config.EnsureDirectory(s.cfg.Server.SocketPath)

		l, err := net.Listen("unix", s.cfg.Server.SocketPath)
		if err != nil {
			return fmt.Errorf("failed to listen on socket %s: %w", s.cfg.Server.SocketPath, err)
		}
		_ = os.Chmod(s.cfg.Server.SocketPath, 0660)
		s.listener = l
		log.Printf("[AGENT] Listening on Unix Domain Socket: %s", s.cfg.Server.SocketPath)
	} else {
		addr := fmt.Sprintf("%s:%d", s.cfg.Server.ListenAddress, s.cfg.Server.Port)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on %s: %w", addr, err)
		}
		s.listener = l
		log.Printf("[AGENT] Listening on TCP: %s", addr)
	}

	return s.httpServer.Serve(s.listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.cfg.Server.UseSocket && s.cfg.Server.SocketPath != "" {
		_ = os.Remove(s.cfg.Server.SocketPath)
	}
	return s.httpServer.Shutdown(ctx)
}

func ErrServerClosed() error {
	return http.ErrServerClosed
}


package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kodepreneur/agent/config"
	"github.com/kodepreneur/agent/server"
)

var (
	version = "1.0.0-dev"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML configuration file")
	isDev := flag.Bool("dev", false, "Run in local development / mock mode")
	port := flag.Int("port", 0, "Override port")
	socket := flag.String("socket", "", "Override socket path")
	showVersion := flag.Bool("version", false, "Print agent version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Kodepreneur Agent v%s\n", version)
		os.Exit(0)
	}

	banner := `
======================================================
  Kodepreneur Privileged Server Agent v%s
  Simple server management without unnecessary complexity
======================================================
`
	fmt.Printf(banner, version)

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("[FATAL] Error loading configuration: %v", err)
	}

	if *isDev {
		cfg.Environment.IsDev = true
	}
	if *port > 0 {
		cfg.Server.Port = *port
		cfg.Server.UseSocket = false
	}
	if *socket != "" {
		cfg.Server.SocketPath = *socket
		cfg.Server.UseSocket = true
	}

	log.Printf("[AGENT] Starting daemon (DevMode=%t)...", cfg.Environment.IsDev)

	srv := server.NewServer(cfg)

	go func() {
		if err := srv.Start(); err != nil && err != server.ErrServerClosed() {
			log.Fatalf("[FATAL] Agent server crashed: %v", err)
		}
	}()

	// Signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[AGENT] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("[AGENT] Agent daemon stopped.")
}

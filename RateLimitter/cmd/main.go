package main

import (
	"log"

	"rate-limiter/internal/config"
	"rate-limiter/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server (this will block until shutdown)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

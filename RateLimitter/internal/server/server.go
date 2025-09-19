package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	config      *config.Config
	storage     storage.Storage
	rateLimiter *limiter.RateLimiter
	router      *gin.Engine
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize storage
	var store storage.Storage
	var err error

	// Try Redis first, fallback to memory storage
	store, err = storage.NewRedisStorage(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Printf("Failed to connect to Redis, using memory storage: %v", err)
		store = storage.NewMemoryStorage()
	}

	// Initialize rate limiter
	rateLimiter := limiter.NewRateLimiter(store, cfg)

	// Initialize router
	router := gin.Default()

	server := &Server{
		config:      cfg,
		storage:     store,
		rateLimiter: rateLimiter,
		router:      router,
	}

	// Setup routes
	server.setupRoutes()

	return server, nil
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() {
	// Health check endpoint (no rate limiting)
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API endpoints with rate limiting
	api := s.router.Group("/api")
	api.Use(middleware.RateLimiterMiddleware(s.rateLimiter))

	// Test endpoint
	api.GET("/test", func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		c.JSON(200, gin.H{
			"message": "Request successful",
			"ip":      ip,
			"token":   token,
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Rate limit status endpoint
	api.GET("/status", func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		remaining, err := s.rateLimiter.GetRemainingRequests(c.Request.Context(), ip, token)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get remaining requests"})
			return
		}

		blocked, err := s.rateLimiter.IsBlocked(c.Request.Context(), ip, token)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to check block status"})
			return
		}

		c.JSON(200, gin.H{
			"remaining_requests": remaining,
			"blocked":            blocked,
			"ip":                 ip,
			"token":              token,
		})
	})

	// Admin endpoint to unblock IPs/tokens (no rate limiting)
	admin := s.router.Group("/admin")
	admin.POST("/unblock", func(c *gin.Context) {
		var request struct {
			Type string `json:"type" binding:"required"` // "ip" or "token"
			Key  string `json:"key" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		ctx := c.Request.Context()
		var err error

		if request.Type == "ip" {
			err = s.storage.Unblock(ctx, request.Key)
		} else if request.Type == "token" {
			err = s.storage.Unblock(ctx, request.Key)
		} else {
			c.JSON(400, gin.H{"error": "Invalid type. Must be 'ip' or 'token'"})
			return
		}

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to unblock"})
			return
		}

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%s %s unblocked successfully", request.Type, request.Key),
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + s.config.Server.Port,
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", s.config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	// Close storage connection
	if err := s.storage.Close(); err != nil {
		log.Printf("Error closing storage: %v", err)
	}

	log.Println("Server exited")
	return nil
}

package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterMiddleware(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IPRequestsPerSecond:    2,
			IPBlockDurationMinutes: 1,
			TokenLimits: map[string]config.TokenLimit{
				"test-token": {
					RequestsPerSecond:    5,
					BlockDurationMinutes: 2,
				},
			},
		},
	}

	// Create memory storage for testing
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	// Create rate limiter
	rl := limiter.NewRateLimiter(storage, cfg)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("IP Rate Limiting", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.RateLimiterMiddleware(rl))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		// First 2 requests should succeed
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
		}

		// 3rd request should be blocked
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code)
		assert.Contains(t, w.Body.String(), "you have reached the maximum number of requests")
	})

	t.Run("Token Rate Limiting", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.RateLimiterMiddleware(rl))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		// First 5 requests with token should succeed
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.2:12345"
			req.Header.Set("API_KEY", "test-token")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
		}

		// 6th request should be blocked
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		req.Header.Set("API_KEY", "test-token")
		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code)
		assert.Contains(t, w.Body.String(), "you have reached the maximum number of requests")
	})

	t.Run("Unknown Token Falls Back to IP", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.RateLimiterMiddleware(rl))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		// First 2 requests with unknown token should succeed (IP limit)
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.3:12345"
			req.Header.Set("API_KEY", "unknown-token")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
		}

		// 3rd request should be blocked (IP limit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.3:12345"
		req.Header.Set("API_KEY", "unknown-token")
		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code)
	})

	t.Run("X-Forwarded-For Header", func(t *testing.T) {
		router := gin.New()
		router.Use(middleware.RateLimiterMiddleware(rl))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})

		// First 2 requests with X-Forwarded-For should succeed
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "10.0.0.1:12345"
			req.Header.Set("X-Forwarded-For", "192.168.1.4")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
		}

		// 3rd request should be blocked
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		req.Header.Set("X-Forwarded-For", "192.168.1.4")
		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code)
	})
}

package test

import (
	"context"
	"testing"
	"time"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/storage"
)

func TestRateLimiter_IPLimit(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IPRequestsPerSecond:    3,
			IPBlockDurationMinutes: 1,
			TokenLimits:            make(map[string]config.TokenLimit),
		},
	}

	// Create memory storage for testing
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	// Create rate limiter
	rl := limiter.NewRateLimiter(storage, cfg)

	ctx := context.Background()
	ip := "192.168.1.1"

	// Test first 3 requests should be allowed
	for i := 0; i < 3; i++ {
		result, err := rl.CheckRequest(ctx, ip, "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !result.Allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 4th request should be blocked
	result, err := rl.CheckRequest(ctx, ip, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Allowed {
		t.Error("4th request should be blocked")
	}

	// Check if IP is blocked
	blocked, err := rl.IsBlocked(ctx, ip, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !blocked {
		t.Error("IP should be blocked")
	}
}

func TestRateLimiter_TokenLimit(t *testing.T) {
	// Create test configuration with token limits
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

	ctx := context.Background()
	ip := "192.168.1.1"
	token := "test-token"

	// Test token limits (5 requests should be allowed)
	for i := 0; i < 5; i++ {
		result, err := rl.CheckRequest(ctx, ip, token)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !result.Allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	result, err := rl.CheckRequest(ctx, ip, token)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Allowed {
		t.Error("6th request should be blocked")
	}

	// Check if token is blocked
	blocked, err := rl.IsBlocked(ctx, ip, token)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !blocked {
		t.Error("Token should be blocked")
	}
}

func TestRateLimiter_UnknownTokenFallsBackToIP(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IPRequestsPerSecond:    2,
			IPBlockDurationMinutes: 1,
			TokenLimits: map[string]config.TokenLimit{
				"known-token": {
					RequestsPerSecond:    10,
					BlockDurationMinutes: 5,
				},
			},
		},
	}

	// Create memory storage for testing
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	// Create rate limiter
	rl := limiter.NewRateLimiter(storage, cfg)

	ctx := context.Background()
	ip := "192.168.1.1"
	unknownToken := "unknown-token"

	// Test that unknown token falls back to IP limits (2 requests)
	for i := 0; i < 2; i++ {
		result, err := rl.CheckRequest(ctx, ip, unknownToken)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !result.Allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 3rd request should be blocked (IP limit)
	result, err := rl.CheckRequest(ctx, ip, unknownToken)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Allowed {
		t.Error("3rd request should be blocked (IP limit)")
	}
}

func TestRateLimiter_GetRemainingRequests(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IPRequestsPerSecond:    5,
			IPBlockDurationMinutes: 1,
			TokenLimits:            make(map[string]config.TokenLimit),
		},
	}

	// Create memory storage for testing
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	// Create rate limiter
	rl := limiter.NewRateLimiter(storage, cfg)

	ctx := context.Background()
	ip := "192.168.1.1"

	// Initially should have 5 remaining requests
	remaining, err := rl.GetRemainingRequests(ctx, ip, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if remaining != 5 {
		t.Errorf("Expected 5 remaining requests, got %d", remaining)
	}

	// Make 2 requests
	for i := 0; i < 2; i++ {
		_, err := rl.CheckRequest(ctx, ip, "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	// Should have 3 remaining requests
	remaining, err = rl.GetRemainingRequests(ctx, ip, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if remaining != 3 {
		t.Errorf("Expected 3 remaining requests, got %d", remaining)
	}
}

func TestRateLimiter_Expiration(t *testing.T) {
	// Create test configuration with very short block duration
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			IPRequestsPerSecond:    2,
			IPBlockDurationMinutes: 0, // No blocking, just rate limiting
			TokenLimits:            make(map[string]config.TokenLimit),
		},
	}

	// Create memory storage for testing
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	// Create rate limiter
	rl := limiter.NewRateLimiter(storage, cfg)

	ctx := context.Background()
	ip := "192.168.1.100" // Use different IP to avoid interference

	// Make 2 requests to reach limit
	for i := 0; i < 2; i++ {
		result, err := rl.CheckRequest(ctx, ip, "")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !result.Allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 3rd request should be blocked
	result, err := rl.CheckRequest(ctx, ip, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Allowed {
		t.Error("3rd request should be blocked")
	}

	// Wait for counter expiration (1 second)
	time.Sleep(2 * time.Second)

	// After counter expiration, requests should be allowed again
	result, err = rl.CheckRequest(ctx, ip, "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !result.Allowed {
		t.Error("Request after counter expiration should be allowed")
	}
}

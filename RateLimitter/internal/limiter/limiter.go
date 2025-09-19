package limiter

import (
	"context"
	"fmt"
	"time"

	"rate-limiter/internal/config"
	"rate-limiter/internal/storage"
)

// LimiterResult represents the result of a rate limit check
type LimiterResult struct {
	Allowed bool
	Reason  string
}

// RateLimiter handles rate limiting logic
type RateLimiter struct {
	storage storage.Storage
	config  *config.Config
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(storage storage.Storage, config *config.Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

// CheckRequest checks if a request should be allowed based on IP and token
func (rl *RateLimiter) CheckRequest(ctx context.Context, ip, token string) (*LimiterResult, error) {
	// First check if IP or token is blocked
	ipBlocked, err := rl.storage.IsBlocked(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("failed to check IP block status: %w", err)
	}

	if ipBlocked {
		return &LimiterResult{
			Allowed: false,
			Reason:  "IP is blocked",
		}, nil
	}

	// Check token block status if token is provided
	if token != "" {
		tokenBlocked, err := rl.storage.IsBlocked(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("failed to check token block status: %w", err)
		}

		if tokenBlocked {
			return &LimiterResult{
				Allowed: false,
				Reason:  "Token is blocked",
			}, nil
		}
	}

	// Determine which limits to apply (token limits override IP limits)
	var requestsPerSecond int
	var blockDurationMinutes int
	var key string

	if token != "" {
		if tokenLimit, exists := rl.config.RateLimit.TokenLimits[token]; exists {
			// Use token-specific limits
			requestsPerSecond = tokenLimit.RequestsPerSecond
			blockDurationMinutes = tokenLimit.BlockDurationMinutes
			key = token
		} else {
			// Use IP limits for unknown tokens
			requestsPerSecond = rl.config.RateLimit.IPRequestsPerSecond
			blockDurationMinutes = rl.config.RateLimit.IPBlockDurationMinutes
			key = ip
		}
	} else {
		// Use IP limits
		requestsPerSecond = rl.config.RateLimit.IPRequestsPerSecond
		blockDurationMinutes = rl.config.RateLimit.IPBlockDurationMinutes
		key = ip
	}

	// Check current request count
	currentCount, err := rl.storage.GetRequestCount(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get request count: %w", err)
	}

	// If limit exceeded, block the key and deny request
	if currentCount >= requestsPerSecond {
		blockDuration := time.Duration(blockDurationMinutes) * time.Minute
		err = rl.storage.Block(ctx, key, blockDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to block key: %w", err)
		}

		return &LimiterResult{
			Allowed: false,
			Reason:  fmt.Sprintf("Rate limit exceeded: %d requests per second", requestsPerSecond),
		}, nil
	}

	// Increment request count
	expiration := time.Second
	err = rl.storage.IncrementRequestCount(ctx, key, expiration)
	if err != nil {
		return nil, fmt.Errorf("failed to increment request count: %w", err)
	}

	return &LimiterResult{
		Allowed: true,
		Reason:  "Request allowed",
	}, nil
}

// GetRemainingRequests returns the number of remaining requests for a key
func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, ip, token string) (int, error) {
	var requestsPerSecond int
	var key string

	if token != "" {
		if tokenLimit, exists := rl.config.RateLimit.TokenLimits[token]; exists {
			requestsPerSecond = tokenLimit.RequestsPerSecond
			key = token
		} else {
			requestsPerSecond = rl.config.RateLimit.IPRequestsPerSecond
			key = ip
		}
	} else {
		requestsPerSecond = rl.config.RateLimit.IPRequestsPerSecond
		key = ip
	}

	currentCount, err := rl.storage.GetRequestCount(ctx, key)
	if err != nil {
		return 0, err
	}

	remaining := requestsPerSecond - currentCount
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// IsBlocked checks if a key is currently blocked
func (rl *RateLimiter) IsBlocked(ctx context.Context, ip, token string) (bool, error) {
	// Check IP block
	ipBlocked, err := rl.storage.IsBlocked(ctx, ip)
	if err != nil {
		return false, err
	}

	if ipBlocked {
		return true, nil
	}

	// Check token block if token is provided
	if token != "" {
		tokenBlocked, err := rl.storage.IsBlocked(ctx, token)
		if err != nil {
			return false, err
		}

		return tokenBlocked, nil
	}

	return false, nil
}

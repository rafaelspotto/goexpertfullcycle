package storage

import (
	"context"
	"time"
)

// Storage defines the interface for rate limiter storage
type Storage interface {
	// GetRequestCount returns the current request count for a key
	GetRequestCount(ctx context.Context, key string) (int, error)

	// IncrementRequestCount increments the request count for a key
	IncrementRequestCount(ctx context.Context, key string, expiration time.Duration) error

	// IsBlocked checks if a key is currently blocked
	IsBlocked(ctx context.Context, key string) (bool, error)

	// Block blocks a key for the specified duration
	Block(ctx context.Context, key string, duration time.Duration) error

	// Unblock removes the block for a key
	Unblock(ctx context.Context, key string) error

	// Close closes the storage connection
	Close() error
}

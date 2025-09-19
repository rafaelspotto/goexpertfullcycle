package storage

import (
	"context"
	"sync"
	"time"
)

// MemoryStorage implements the Storage interface using in-memory storage
// This is useful for testing or when Redis is not available
type MemoryStorage struct {
	mu          sync.RWMutex
	counters    map[string]int
	blocks      map[string]time.Time
	expirations map[string]time.Time
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		counters:    make(map[string]int),
		blocks:      make(map[string]time.Time),
		expirations: make(map[string]time.Time),
	}

	// Start cleanup goroutine
	go storage.cleanup()

	return storage
}

// GetRequestCount returns the current request count for a key
func (m *MemoryStorage) GetRequestCount(ctx context.Context, key string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if key has expired
	if exp, exists := m.expirations[key]; exists && time.Now().After(exp) {
		return 0, nil
	}

	return m.counters[key], nil
}

// IncrementRequestCount increments the request count for a key
func (m *MemoryStorage) IncrementRequestCount(ctx context.Context, key string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if key has expired
	if exp, exists := m.expirations[key]; exists && time.Now().After(exp) {
		m.counters[key] = 0
	}

	m.counters[key]++
	m.expirations[key] = time.Now().Add(expiration)

	return nil
}

// IsBlocked checks if a key is currently blocked
func (m *MemoryStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	blockTime, exists := m.blocks[key]
	if !exists {
		return false, nil
	}

	// Check if block has expired
	if time.Now().After(blockTime) {
		return false, nil
	}

	return true, nil
}

// Block blocks a key for the specified duration
func (m *MemoryStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.blocks[key] = time.Now().Add(duration)
	return nil
}

// Unblock removes the block for a key
func (m *MemoryStorage) Unblock(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.blocks, key)
	return nil
}

// Close closes the storage (no-op for memory storage)
func (m *MemoryStorage) Close() error {
	return nil
}

// cleanup removes expired entries periodically
func (m *MemoryStorage) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()

		// Clean up expired counters
		for key, exp := range m.expirations {
			if now.After(exp) {
				delete(m.counters, key)
				delete(m.expirations, key)
			}
		}

		// Clean up expired blocks
		for key, blockTime := range m.blocks {
			if now.After(blockTime) {
				delete(m.blocks, key)
			}
		}

		m.mu.Unlock()
	}
}

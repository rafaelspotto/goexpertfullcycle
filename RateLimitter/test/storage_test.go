package test

import (
	"context"
	"testing"
	"time"

	"rate-limiter/internal/storage"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()
	key := "test-key"

	t.Run("GetRequestCount - Non-existent key", func(t *testing.T) {
		count, err := storage.GetRequestCount(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("IncrementRequestCount", func(t *testing.T) {
		err := storage.IncrementRequestCount(ctx, key, time.Second)
		assert.NoError(t, err)

		count, err := storage.GetRequestCount(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Multiple Increments", func(t *testing.T) {
		for i := 0; i < 4; i++ {
			err := storage.IncrementRequestCount(ctx, key, time.Second)
			assert.NoError(t, err)
		}

		count, err := storage.GetRequestCount(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("IsBlocked - Non-blocked key", func(t *testing.T) {
		blocked, err := storage.IsBlocked(ctx, key)
		assert.NoError(t, err)
		assert.False(t, blocked)
	})

	t.Run("Block key", func(t *testing.T) {
		err := storage.Block(ctx, key, time.Minute)
		assert.NoError(t, err)

		blocked, err := storage.IsBlocked(ctx, key)
		assert.NoError(t, err)
		assert.True(t, blocked)
	})

	t.Run("Unblock key", func(t *testing.T) {
		err := storage.Unblock(ctx, key)
		assert.NoError(t, err)

		blocked, err := storage.IsBlocked(ctx, key)
		assert.NoError(t, err)
		assert.False(t, blocked)
	})

	t.Run("Expiration", func(t *testing.T) {
		// Use a different key to avoid interference from previous tests
		expKey := "expiration-test-key"

		// Set a very short expiration
		err := storage.IncrementRequestCount(ctx, expKey, 100*time.Millisecond)
		assert.NoError(t, err)

		// Should have count immediately
		count, err := storage.GetRequestCount(ctx, expKey)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		// Wait for expiration
		time.Sleep(200 * time.Millisecond)

		// Count should be 0 after expiration
		count, err = storage.GetRequestCount(ctx, expKey)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestMemoryStorage_Concurrent(t *testing.T) {
	storage := storage.NewMemoryStorage()
	defer storage.Close()

	ctx := context.Background()
	key := "concurrent-test"

	// Test concurrent increments
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			err := storage.IncrementRequestCount(ctx, key, time.Second)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Check final count
	count, err := storage.GetRequestCount(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, 10, count)
}

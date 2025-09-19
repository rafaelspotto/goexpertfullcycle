package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStorage implements the Storage interface using Redis
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage creates a new Redis storage instance
func NewRedisStorage(host, port, password string, db int) (*RedisStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStorage{client: rdb}, nil
}

// GetRequestCount returns the current request count for a key
func (r *RedisStorage) GetRequestCount(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid count value: %w", err)
	}

	return count, nil
}

// IncrementRequestCount increments the request count for a key
func (r *RedisStorage) IncrementRequestCount(ctx context.Context, key string, expiration time.Duration) error {
	pipe := r.client.Pipeline()

	// Increment counter
	pipe.Incr(ctx, key)
	// Set expiration only if key doesn't exist
	pipe.Expire(ctx, key, expiration)

	_, err := pipe.Exec(ctx)
	return err
}

// IsBlocked checks if a key is currently blocked
func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := fmt.Sprintf("block:%s", key)
	exists, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

// Block blocks a key for the specified duration
func (r *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	blockKey := fmt.Sprintf("block:%s", key)
	return r.client.Set(ctx, blockKey, "1", duration).Err()
}

// Unblock removes the block for a key
func (r *RedisStorage) Unblock(ctx context.Context, key string) error {
	blockKey := fmt.Sprintf("block:%s", key)
	return r.client.Del(ctx, blockKey).Err()
}

// Close closes the Redis connection
func (r *RedisStorage) Close() error {
	return r.client.Close()
}

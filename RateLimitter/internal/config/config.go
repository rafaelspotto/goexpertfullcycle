package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the rate limiter
type Config struct {
	Redis     RedisConfig
	RateLimit RateLimitConfig
	Server    ServerConfig
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	IPRequestsPerSecond    int
	IPBlockDurationMinutes int
	TokenLimits            map[string]TokenLimit
}

// TokenLimit holds configuration for a specific token
type TokenLimit struct {
	RequestsPerSecond    int
	BlockDurationMinutes int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file if it exists
	_ = godotenv.Load("config.env")

	config := &Config{
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		RateLimit: RateLimitConfig{
			IPRequestsPerSecond:    getEnvAsInt("RATE_LIMIT_IP_REQUESTS_PER_SECOND", 5),
			IPBlockDurationMinutes: getEnvAsInt("RATE_LIMIT_IP_BLOCK_DURATION_MINUTES", 5),
			TokenLimits:            loadTokenLimits(),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}

	return config, nil
}

// loadTokenLimits loads token-specific rate limits from environment variables
func loadTokenLimits() map[string]TokenLimit {
	tokenLimits := make(map[string]TokenLimit)

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TOKEN_LIMIT_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimPrefix(parts[0], "TOKEN_LIMIT_")
			value := parts[1]

			// Parse format: REQUESTS_PER_SECOND:BLOCK_DURATION_MINUTES
			limitParts := strings.Split(value, ":")
			if len(limitParts) != 2 {
				continue
			}

			requestsPerSecond, err := strconv.Atoi(limitParts[0])
			if err != nil {
				continue
			}

			blockDurationMinutes, err := strconv.Atoi(limitParts[1])
			if err != nil {
				continue
			}

			tokenLimits[key] = TokenLimit{
				RequestsPerSecond:    requestsPerSecond,
				BlockDurationMinutes: blockDurationMinutes,
			}
		}
	}

	return tokenLimits
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

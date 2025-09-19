#!/bin/bash

# Test script for rate limiter

echo "Running unit tests..."
go test -v ./test/...

echo ""
echo "Running integration tests..."

# Start Redis in background
docker-compose up -d redis

# Wait for Redis to be ready
echo "Waiting for Redis to be ready..."
until docker-compose exec redis redis-cli ping; do
  sleep 1
done

# Run integration tests
go test -v -tags=integration ./test/...

# Cleanup
docker-compose down

echo "All tests completed!"

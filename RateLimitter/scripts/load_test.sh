#!/bin/bash

# Load testing script for rate limiter

echo "Starting load test..."

# Start the application
docker-compose up -d

# Wait for application to be ready
echo "Waiting for application to be ready..."
until curl -s http://localhost:8080/health > /dev/null; do
  sleep 1
done

echo "Application is ready. Starting load tests..."

# Test IP rate limiting (5 requests per second)
echo "Testing IP rate limiting (5 req/s)..."
for i in {1..10}; do
  response=$(curl -s -w "%{http_code}" -o /dev/null http://localhost:8080/api/test)
  echo "Request $i: HTTP $response"
  sleep 0.1
done

echo ""
echo "Testing token rate limiting (10 req/s for abc123)..."
for i in {1..15}; do
  response=$(curl -s -w "%{http_code}" -o /dev/null -H "API_KEY: abc123" http://localhost:8080/api/test)
  echo "Request $i: HTTP $response"
  sleep 0.1
done

echo ""
echo "Testing unknown token (falls back to IP limit)..."
for i in {1..10}; do
  response=$(curl -s -w "%{http_code}" -o /dev/null -H "API_KEY: unknown-token" http://localhost:8080/api/test)
  echo "Request $i: HTTP $response"
  sleep 0.1
done

echo ""
echo "Testing rate limit status endpoint..."
curl -s http://localhost:8080/api/status | jq .

echo ""
echo "Testing admin unblock endpoint..."
curl -s -X POST http://localhost:8080/admin/unblock \
  -H "Content-Type: application/json" \
  -d '{"type": "ip", "key": "127.0.0.1"}' | jq .

echo ""
echo "Load test completed!"

# Cleanup
docker-compose down

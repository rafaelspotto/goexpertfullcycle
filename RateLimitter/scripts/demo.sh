#!/bin/bash

# Demo script for rate limiter

echo "🚀 Rate Limiter Demo"
echo "===================="
echo ""

# Start the application
echo "Starting application with Docker Compose..."
docker-compose up -d

# Wait for application to be ready
echo "Waiting for application to be ready..."
until curl -s http://localhost:8080/health > /dev/null; do
  sleep 1
done

echo "✅ Application is ready!"
echo ""

# Test 1: IP Rate Limiting
echo "📊 Test 1: IP Rate Limiting (5 req/s)"
echo "Making 7 requests to test IP rate limiting..."
for i in {1..7}; do
  response=$(curl -s -w "Request $i: HTTP %{http_code}" -o /dev/null http://localhost:8080/api/test)
  echo "  $response"
  sleep 0.1
done
echo ""

# Test 2: Token Rate Limiting
echo "📊 Test 2: Token Rate Limiting (10 req/s for abc123)"
echo "Making 12 requests with token abc123..."
for i in {1..12}; do
  response=$(curl -s -w "Request $i: HTTP %{http_code}" -o /dev/null -H "API_KEY: abc123" http://localhost:8080/api/test)
  echo "  $response"
  sleep 0.1
done
echo ""

# Test 3: Unknown Token (falls back to IP)
echo "📊 Test 3: Unknown Token (falls back to IP limit)"
echo "Making 7 requests with unknown token..."
for i in {1..7}; do
  response=$(curl -s -w "Request $i: HTTP %{http_code}" -o /dev/null -H "API_KEY: unknown-token" http://localhost:8080/api/test)
  echo "  $response"
  sleep 0.1
done
echo ""

# Test 4: Rate Limit Status
echo "📊 Test 4: Rate Limit Status"
echo "Checking rate limit status..."
curl -s http://localhost:8080/api/status | jq .
echo ""

# Test 5: Admin Unblock
echo "📊 Test 5: Admin Unblock"
echo "Unblocking IP 127.0.0.1..."
curl -s -X POST http://localhost:8080/admin/unblock \
  -H "Content-Type: application/json" \
  -d '{"type": "ip", "key": "127.0.0.1"}' | jq .
echo ""

# Test 6: Health Check
echo "📊 Test 6: Health Check"
echo "Checking application health..."
curl -s http://localhost:8080/health | jq .
echo ""

echo "🎉 Demo completed!"
echo ""
echo "To stop the application, run: docker-compose down"
echo "To view logs, run: docker-compose logs -f"

.PHONY: help build run test test-integration clean docker-build docker-up docker-down load-test

# Default target
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application locally"
	@echo "  test            - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-up       - Start services with Docker Compose"
	@echo "  docker-down     - Stop services"
	@echo "  load-test       - Run load tests"
	@echo "  all-tests       - Run all tests"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/rate-limiter cmd/main.go

# Run the application locally
run:
	@echo "Running application..."
	go run cmd/main.go

# Run unit tests
test:
	@echo "Running unit tests..."
	go test -v ./test/...

# Run integration tests
test-integration:
	@echo "Starting Redis for integration tests..."
	docker-compose up -d redis
	@echo "Waiting for Redis to be ready..."
	@until docker-compose exec redis redis-cli ping > /dev/null 2>&1; do sleep 1; done
	@echo "Running integration tests..."
	go test -v -tags=integration ./test/...
	@echo "Stopping Redis..."
	docker-compose down

# Run all tests
all-tests: test test-integration

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t rate-limiter .

# Start services with Docker Compose
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# Stop services
docker-down:
	@echo "Stopping services..."
	docker-compose down

# Run load tests
load-test:
	@echo "Running load tests..."
	./scripts/load_test.sh

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Full test suite
full-test: fmt lint test test-integration

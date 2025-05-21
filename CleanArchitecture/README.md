# Clean Architecture - Order Management System

This project implements a Clean Architecture system for order management with multiple interfaces:
- REST API (Port 8080)
- gRPC Service (Port 50051)
- GraphQL API (Port 8081)

## Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose

## Getting Started

1. Start the database:
```bash
docker compose up -d
```

2. Run database migrations:
```bash
go run cmd/migrate/main.go
```

3. Start the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

### REST API (Port 8080)
- GET /order - List all orders

### gRPC (Port 50051)
- ListOrders - List all orders

### GraphQL (Port 8081)
- Query: listOrders - List all orders

## Project Structure
```
.
├── cmd/
│   ├── api/
│   └── migrate/
├── internal/
│   ├── core/
│   │   ├── domain/
│   │   ├── usecase/
│   │   └── repository/
│   ├── infrastructure/
│   │   ├── database/
│   │   └── grpc/
│   └── interfaces/
│       ├── http/
│       ├── grpc/
│       └── graphql/
├── pkg/
├── docker-compose.yaml
└── go.mod
```
.PHONY: generate run docker-up docker-down

generate:
	@echo "Generating gRPC code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/order.proto
	@echo "Generating GraphQL code..."
	@gqlgen generate

run:
	@go run cmd/api/main.go

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

install-tools:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/99designs/gqlgen@latest 
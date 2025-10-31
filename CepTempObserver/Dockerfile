FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE=server
RUN if [ "$SERVICE" = "service-a" ]; then \
      go build -o main cmd/service-a/main.go; \
    else \
      go build -o main cmd/server/main.go; \
    fi

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
EXPOSE 8081
CMD ["./main"]


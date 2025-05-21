package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/goexpertfullcycle/cleanarchitecture/graph"
	"github.com/goexpertfullcycle/cleanarchitecture/internal/core/usecase"
	"github.com/goexpertfullcycle/cleanarchitecture/internal/infrastructure/database"
	"github.com/goexpertfullcycle/cleanarchitecture/internal/interfaces/graphql"
	"github.com/goexpertfullcycle/cleanarchitecture/internal/interfaces/grpc"
	httpHandler "github.com/goexpertfullcycle/cleanarchitecture/internal/interfaces/http"
)

func main() {
	// Initialize database
	repo, err := database.NewPostgresRepository()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize use case
	orderUseCase := usecase.NewOrderUseCase(repo)

	// Initialize HTTP handler
	orderHandler := httpHandler.NewOrderHandler(orderUseCase)

	// Setup HTTP routes
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderHandler.List(w, r)
		case http.MethodPost:
			orderHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Initialize GraphQL
	resolver := graphql.NewResolver(orderUseCase)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Start gRPC server in a goroutine
	go func() {
		if err := grpc.StartGRPCServer(orderUseCase); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Start HTTP server
	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

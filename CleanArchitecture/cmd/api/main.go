package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/core/usecase"
	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/infrastructure/database"
	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/interfaces/grpc"
	httpHandler "github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/interfaces/http"
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

	// Initialize GraphQL (simple implementation)
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Simple GraphQL query handling
		if request.Query == "{ orders { id price tax finalPrice createdAt updatedAt } }" {
			orders, err := orderUseCase.List()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"data": map[string]interface{}{
					"orders": orders,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Handle createOrder mutation
		if request.Query == "mutation { createOrder(input: {price: 100.0, tax: 10.0}) { id price tax finalPrice createdAt updatedAt } }" {
			order, err := orderUseCase.Create(100.0, 10.0)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"data": map[string]interface{}{
					"createOrder": order,
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Default response
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"message": "GraphQL endpoint is working! Try: { orders { id price tax finalPrice createdAt updatedAt } }",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

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

package main

import (
	"ceptemperature/internal/handlers"
	"ceptemperature/internal/telemetry"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Inicializar tracing OpenTelemetry
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "otel-collector:4318"
	}

	tp, err := telemetry.InitTracing("service-a", otlpEndpoint)
	if err != nil {
		log.Printf("Failed to initialize tracing: %v", err)
	} else {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	// Configurar Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Middleware de tracing OpenTelemetry (deve ser o primeiro)
	r.Use(otelgin.Middleware("service-a"))

	// Middleware de logging
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Middleware de CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	r.Use(cors.New(config))

	// Registrar rotas - POST endpoint para receber CEP
	r.POST("/", handlers.HandleCepRequest)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Obter porta das variáveis de ambiente
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Criar canal para capturar sinais de interrupção
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Service A starting on port %s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-sigChan
	log.Println("Shutting down Service A...")

	// Shutdown do tracer provider
	if tp != nil {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}


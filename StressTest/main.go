package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Result representa o resultado de uma requisição HTTP
type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

// Stats contém as estatísticas coletadas durante o teste
type Stats struct {
	TotalRequests     int
	Successful200     int
	StatusCodes       map[int]int
	TotalDuration     time.Duration
	MinDuration       time.Duration
	MaxDuration       time.Duration
	AvgDuration       time.Duration
	Errors            int
	RequestsPerSecond float64
}

func main() {
	// Parse dos argumentos da linha de comando
	var url string
	var requests int
	var concurrency int

	flag.StringVar(&url, "url", "", "URL do serviço a ser testado")
	flag.IntVar(&requests, "requests", 0, "Número total de requests")
	flag.IntVar(&concurrency, "concurrency", 0, "Número de chamadas simultâneas")
	flag.Parse()

	// Validação dos parâmetros
	if url == "" || requests <= 0 || concurrency <= 0 {
		fmt.Println("Uso: stresstest --url=<URL> --requests=<N> --concurrency=<N>")
		fmt.Println("Exemplo: stresstest --url=http://google.com --requests=1000 --concurrency=10")
		flag.PrintDefaults()
		return
	}

	// Execução do teste de carga
	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Requests: %d\n", requests)
	fmt.Printf("Concorrência: %d\n", concurrency)
	fmt.Println(strings.Repeat("=", 50))

	stats := runLoadTest(url, requests, concurrency)

	// Geração do relatório
	printReport(stats)
}

// runLoadTest executa o teste de carga
func runLoadTest(url string, totalRequests, concurrency int) *Stats {
	startTime := time.Now()

	// Canal para receber resultados
	results := make(chan Result, totalRequests)

	// WaitGroup para sincronizar as goroutines
	var wg sync.WaitGroup

	// Canal para controlar a concorrência
	semaphore := make(chan struct{}, concurrency)

	// Inicia as goroutines para fazer as requisições
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Adquire o semáforo
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Executa a requisição
			result := makeRequest(url)
			results <- result
		}()
	}

	// Aguarda todas as goroutines terminarem
	go func() {
		wg.Wait()
		close(results)
	}()

	// Coleta os resultados
	stats := &Stats{
		StatusCodes: make(map[int]int),
		MinDuration: time.Hour, // Inicializa com um valor alto
	}

	var totalDuration time.Duration
	firstDuration := true

	for result := range results {
		stats.TotalRequests++

		if result.Error != nil {
			stats.Errors++
		} else {
			stats.StatusCodes[result.StatusCode]++
			if result.StatusCode == 200 {
				stats.Successful200++
			}

			// Atualiza estatísticas de duração
			if firstDuration {
				stats.MinDuration = result.Duration
				stats.MaxDuration = result.Duration
				firstDuration = false
			} else {
				if result.Duration < stats.MinDuration {
					stats.MinDuration = result.Duration
				}
				if result.Duration > stats.MaxDuration {
					stats.MaxDuration = result.Duration
				}
			}

			totalDuration += result.Duration
		}
	}

	stats.TotalDuration = time.Since(startTime)

	// Calcula duração média
	successfulRequests := stats.TotalRequests - stats.Errors
	if successfulRequests > 0 {
		stats.AvgDuration = totalDuration / time.Duration(successfulRequests)
	}

	// Calcula requests por segundo
	if stats.TotalDuration > 0 {
		stats.RequestsPerSecond = float64(stats.TotalRequests) / stats.TotalDuration.Seconds()
	}

	return stats
}

// makeRequest executa uma requisição HTTP
func makeRequest(url string) Result {
	start := time.Now()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{
			StatusCode: 0,
			Duration:   duration,
			Error:      err,
		}
	}
	defer resp.Body.Close()

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}

// printReport exibe o relatório final
func printReport(stats *Stats) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("RELATÓRIO DO TESTE DE CARGA")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("Tempo total de execução: %v\n", stats.TotalDuration)
	fmt.Printf("Total de requests realizados: %d\n", stats.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", stats.Successful200)
	fmt.Printf("Taxa de sucesso (200): %.2f%%\n", float64(stats.Successful200)/float64(stats.TotalRequests)*100)
	fmt.Printf("Requests por segundo: %.2f\n", stats.RequestsPerSecond)

	if stats.Errors > 0 {
		fmt.Printf("Erros: %d\n", stats.Errors)
	}

	fmt.Println("\nDistribuição de códigos de status HTTP:")
	for statusCode, count := range stats.StatusCodes {
		percentage := float64(count) / float64(stats.TotalRequests) * 100
		fmt.Printf("  %d: %d (%.2f%%)\n", statusCode, count, percentage)
	}

	if stats.TotalRequests > stats.Errors {
		fmt.Printf("\nEstatísticas de duração:\n")
		fmt.Printf("  Mínima: %v\n", stats.MinDuration)
		fmt.Printf("  Máxima: %v\n", stats.MaxDuration)
		fmt.Printf("  Média: %v\n", stats.AvgDuration)
	}

	fmt.Println(strings.Repeat("=", 50))
}

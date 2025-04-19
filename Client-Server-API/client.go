package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type BidResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar request:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Erro ao fazer request:", err)
	}
	defer resp.Body.Close()

	var result BidResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("Erro ao decodificar resposta:", err)
	}

	// Escrevendo no arquivo
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal("Erro ao criar arquivo:", err)
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: " + result.Bid)
	if err != nil {
		log.Fatal("Erro ao escrever no arquivo:", err)
	}

	log.Println("Cotação salva com sucesso:", result.Bid)
}

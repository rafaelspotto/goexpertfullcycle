package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CotacaoResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", handler)
	log.Println("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

	// Timeout de 200ms para a chamada externa
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, "Erro ao criar request", http.StatusInternalServerError)
		log.Println("Erro ao criar request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
		log.Println("Erro na requisição externa:", err)
		return
	}
	defer resp.Body.Close()

	var cotacao CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		http.Error(w, "Erro ao decodificar resposta", http.StatusInternalServerError)
		log.Println("Erro ao decodificar JSON:", err)
		return
	}

	bid := cotacao.USDBRL.Bid

	// Salva a cotação no banco de dados
	if err := salvarCotacaoNoBanco(bid); err != nil {
		log.Println("Erro ao salvar cotação no banco:", err)
	}

	result := map[string]string{
		"bid": bid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func salvarCotacaoNoBanco(bid string) error {
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Cria a tabela se não existir
	createTable := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Insere a cotação com timeout
	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacoes(bid) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()


	_, err = stmt.Exec(bid)
	return err
}

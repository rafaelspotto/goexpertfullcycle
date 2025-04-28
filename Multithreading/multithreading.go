package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Endereco struct {
	Cep string `json:"cep"`
}

type Resultado struct {
	API      string
	Endereco Endereco
	Err      error
}

func buscarComContexto(ctx context.Context, apiNome, url string, ch chan Resultado) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- Resultado{API: apiNome, Err: err}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- Resultado{API: apiNome, Err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- Resultado{API: apiNome, Err: err}
		return
	}

	var endereco Endereco
	err = json.Unmarshal(body, &endereco)
	if err != nil {
		ch <- Resultado{API: apiNome, Err: err}
		return
	}

	ch <- Resultado{API: apiNome, Endereco: endereco}
}

func main() {
	cep := "72110230"
	ch := make(chan Resultado, 2) // buffer 2 para evitar goroutines bloqueadas

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	brasilApiURL := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	viaCepURL := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	go buscarComContexto(ctx, "BrasilAPI", brasilApiURL, ch)
	go buscarComContexto(ctx, "ViaCEP", viaCepURL, ch)

	select {
	case resultado := <-ch:
		if resultado.Err != nil {
			fmt.Println("Erro:", resultado.Err)
			return
		}
		fmt.Println("Resposta recebida da API:", resultado.API)
		fmt.Printf("Endereço: %+v\n", resultado.Endereco)
	case <-ctx.Done():
		fmt.Println("Timeout: Nenhuma resposta em 1 segundo.")
	}
}

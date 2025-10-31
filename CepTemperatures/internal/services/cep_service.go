package services

import (
	"ceptemperature/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//- Função para consultar ViaCEP API

func GetCep(cep string) (*models.ViaCEPResponse, error) {
	isValid, _ := ValidateCep(cep)
	if !isValid {
		return nil, fmt.Errorf("cep must be 8 digits")
	}

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := "https://viacep.com.br/ws/" + cep + "/json/"
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CEP data: %w", err)
	}
	defer resp.Body.Close()

	// Verificar status da resposta
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("viacep API returned status %d", resp.StatusCode)
	}

	var cepResponse models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&cepResponse); err != nil {
		return nil, fmt.Errorf("failed to decode CEP response: %w", err)
	}

	return &cepResponse, nil
}

func ValidateCep(cep string) (bool, error) {
	if len(cep) != 8 {
		return false, fmt.Errorf("cep must be 8 digits")
	}
	return true, nil
}

//Função para validar formato do CEP (8 dígitos)

func ValidateCepFormat(cep string) (bool, error) {
	if len(cep) != 8 {
		return false, fmt.Errorf("cep must be 8 digits")
	}
	return true, nil
}

//Tratamento de erros específicos

func HandleError(err error) (models.ErrorResponse, error) {
	return models.ErrorResponse{Error: err.Error()}, err
}

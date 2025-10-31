package services

import (
	"ceptemperature/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

//- Função para consultar ViaCEP API

func GetCep(cep string) (*models.ViaCEPResponse, error) {
	isValid, err := ValidateCep(cep)
	if !isValid {
		return nil, fmt.Errorf("cep must be 8 digits")
	}

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cepResponse models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&cepResponse); err != nil {
		return nil, err
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

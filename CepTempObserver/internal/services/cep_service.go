package services

import (
	"ceptemperature/internal/models"
	"ceptemperature/internal/telemetry"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

//- Função para consultar ViaCEP API

func GetCep(ctx context.Context, cep string) (*models.ViaCEPResponse, error) {
	// Criar span para medir tempo de resposta da busca de CEP
	ctx, span := telemetry.Tracer.Start(ctx, "cep-service.get-cep")
	defer span.End()

	span.SetAttributes(
		attribute.String("cep", cep),
		attribute.String("service.name", "viacep-api"),
	)

	isValid, _ := ValidateCep(cep)
	if !isValid {
		span.RecordError(fmt.Errorf("cep must be 8 digits"))
		span.SetStatus(codes.Error, "Invalid CEP format")
		return nil, fmt.Errorf("cep must be 8 digits")
	}

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := "https://viacep.com.br/ws/" + cep + "/json/"
	span.SetAttributes(attribute.String("http.url", url))
	
	startTime := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(startTime)
	
	span.SetAttributes(
		attribute.Int64("http.request.duration_ms", duration.Milliseconds()),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch CEP data")
		return nil, fmt.Errorf("failed to fetch CEP data: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	// Verificar status da resposta
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("viacep API returned status %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Non-OK HTTP status")
		return nil, err
	}

	var cepResponse models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&cepResponse); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to decode response")
		return nil, fmt.Errorf("failed to decode CEP response: %w", err)
	}

	span.SetStatus(codes.Ok, "CEP data retrieved successfully")
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

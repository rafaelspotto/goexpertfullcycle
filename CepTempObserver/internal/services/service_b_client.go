package services

import (
	"ceptemperature/internal/models"
	"ceptemperature/internal/telemetry"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
)

// ForwardToServiceB encaminha a requisição para o Serviço B com propagação de tracing
func ForwardToServiceB(ctx context.Context, cep string) (*models.WeatherResponse, error) {
	// Criar span para medir tempo de resposta da chamada ao Serviço B
	ctx, span := telemetry.Tracer.Start(ctx, "service-a.forward-to-service-b")
	defer span.End()

	span.SetAttributes(
		attribute.String("cep", cep),
		attribute.String("service.name", "service-b"),
	)

	// Obter URL do Serviço B das variáveis de ambiente
	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://service-b:8080"
	}

	// Construir URL do endpoint
	url := serviceBURL + "/weather/" + cep
	span.SetAttributes(attribute.String("http.url", url))

	// Criar requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create HTTP request")
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Propagar contexto de tracing para o Serviço B
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	startTime := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	span.SetAttributes(
		attribute.Int64("http.request.duration_ms", duration.Milliseconds()),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to call Service B")
		return nil, fmt.Errorf("failed to call Service B: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	// Tratar diferentes códigos de resposta
	if resp.StatusCode == http.StatusUnprocessableEntity {
		var errorResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			err := fmt.Errorf("invalid zipcode")
			span.RecordError(err)
			span.SetStatus(codes.Error, "Invalid zipcode from Service B")
			return nil, err
		}
		err := fmt.Errorf("invalid zipcode")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid zipcode")
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		var errorResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			err := fmt.Errorf("can not find zipcode")
			span.RecordError(err)
			span.SetStatus(codes.Error, "Zipcode not found")
			return nil, err
		}
		err := fmt.Errorf("can not find zipcode")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Zipcode not found")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Service B returned status %d", resp.StatusCode)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Non-OK HTTP status from Service B")
		return nil, err
	}

	// Decodificar resposta de sucesso
	var weatherResponse models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to decode Service B response")
		return nil, fmt.Errorf("failed to decode Service B response: %w", err)
	}

	span.SetAttributes(
		attribute.String("city", weatherResponse.City),
		attribute.Float64("temperature.celsius", weatherResponse.TempC),
	)

	span.SetStatus(codes.Ok, "Successfully forwarded to Service B")
	return &weatherResponse, nil
}


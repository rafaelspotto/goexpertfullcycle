package services

import (
	"ceptemperature/internal/models"
	"ceptemperature/internal/telemetry"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func GetWeather(ctx context.Context, cep string) (*models.WeatherResponse, error) {
	// Criar span para o serviço de weather
	ctx, span := telemetry.Tracer.Start(ctx, "weather-service.get-weather")
	defer span.End()

	span.SetAttributes(
		attribute.String("cep", cep),
	)

	cepResponse, err := GetCep(ctx, cep)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Verificar se o CEP foi encontrado
	if cepResponse.Erro || cepResponse.Localidade == "" {
		err := fmt.Errorf("cep not found")
		span.RecordError(err)
		span.SetStatus(codes.Error, "CEP not found")
		return nil, err
	}

	span.SetAttributes(attribute.String("localidade", cepResponse.Localidade))

	weatherAPIResponse, err := GetWeatherAPI(ctx, cepResponse.Localidade)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Extrair temperatura em Celsius e converter
	tempC := weatherAPIResponse.Current.TempC
	tempF, tempK, _ := ConvertTemperature(tempC)

	span.SetAttributes(
		attribute.Float64("temperature.celsius", tempC),
		attribute.Float64("temperature.fahrenheit", tempF),
		attribute.Float64("temperature.kelvin", tempK),
	)

	weatherResponse := &models.WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	span.SetStatus(codes.Ok, "Weather data retrieved successfully")
	return weatherResponse, nil
}

func GetWeatherAPI(ctx context.Context, localidade string) (*models.WeatherAPIResponse, error) {
	// Criar span para medir tempo de resposta da busca de temperatura
	ctx, span := telemetry.Tracer.Start(ctx, "weather-service.get-weather-api")
	defer span.End()

	span.SetAttributes(
		attribute.String("localidade", localidade),
		attribute.String("service.name", "weatherapi"),
	)

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		err := fmt.Errorf("WEATHER_API_KEY environment variable is not set")
		span.RecordError(err)
		span.SetStatus(codes.Error, "API key not configured")
		return nil, err
	}

	// Remover espaços em branco e quebras de linha da chave (caso tenha sido copiada com espaços)
	apiKey = strings.TrimSpace(apiKey)

	// Codificar a localidade na URL
	encodedLocalidade := url.QueryEscape(localidade)
	apiUrl := "https://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + encodedLocalidade

	span.SetAttributes(attribute.String("http.url", "https://api.weatherapi.com/v1/current.json"))

	startTime := time.Now()
	resp, err := client.Get(apiUrl)
	duration := time.Since(startTime)

	span.SetAttributes(
		attribute.Int64("http.request.duration_ms", duration.Milliseconds()),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch weather data")
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	// Verificar se a resposta foi bem-sucedida
	if resp.StatusCode != http.StatusOK {
		// Tentar ler o corpo da resposta de erro para mais detalhes
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		var err error
		if errorMsg, ok := errorBody["error"].(map[string]interface{}); ok {
			if msg, ok := errorMsg["message"].(string); ok {
				err = fmt.Errorf("weather API error: %s (status %d)", msg, resp.StatusCode)
			}
		}
		if err == nil {
			err = fmt.Errorf("weather API returned status %d", resp.StatusCode)
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "Non-OK HTTP status")
		return nil, err
	}

	var weatherAPIResponse models.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherAPIResponse); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to decode response")
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	span.SetStatus(codes.Ok, "Weather data retrieved successfully")
	return &weatherAPIResponse, nil
}

func ConvertTemperature(tempC float64) (float64, float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15
	return tempF, tempK, tempC
}

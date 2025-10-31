package services

import (
	"ceptemperature/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetWeather(cep string) (*models.WeatherResponse, error) {
	cepResponse, err := GetCep(cep)
	if err != nil {
		return nil, err
	}

	// Verificar se o CEP foi encontrado
	if cepResponse.Erro || cepResponse.Localidade == "" {
		return nil, fmt.Errorf("cep not found")
	}

	weatherAPIResponse, err := GetWeatherAPI(cepResponse.Localidade)
	if err != nil {
		return nil, err
	}

	// Extrair temperatura em Celsius e converter
	tempC := weatherAPIResponse.Current.TempC
	tempF, tempK, _ := ConvertTemperature(tempC)

	weatherResponse := &models.WeatherResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	return weatherResponse, nil
}

func GetWeatherAPI(localidade string) (*models.WeatherAPIResponse, error) {
	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY environment variable is not set")
	}

	// Remover espaços em branco e quebras de linha da chave (caso tenha sido copiada com espaços)
	apiKey = strings.TrimSpace(apiKey)

	// Codificar a localidade na URL
	encodedLocalidade := url.QueryEscape(localidade)
	apiUrl := "https://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + encodedLocalidade
	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	// Verificar se a resposta foi bem-sucedida
	if resp.StatusCode != http.StatusOK {
		// Tentar ler o corpo da resposta de erro para mais detalhes
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		if errorMsg, ok := errorBody["error"].(map[string]interface{}); ok {
			if msg, ok := errorMsg["message"].(string); ok {
				return nil, fmt.Errorf("weather API error: %s (status %d)", msg, resp.StatusCode)
			}
		}
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var weatherAPIResponse models.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherAPIResponse); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return &weatherAPIResponse, nil
}

func ConvertTemperature(tempC float64) (float64, float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15
	return tempF, tempK, tempC
}

package services

import (
	"ceptemperature/internal/models"
	"encoding/json"
	"net/http"
	"os"
)

func GetWeather(cep string) (*models.WeatherResponse, error) {
	cepResponse, err := GetCep(cep)
	if err != nil {
		return nil, err
	}

	weatherResponse, err := GetWeatherAPI(cepResponse.Localidade)
	if err != nil {
		return nil, err
	}

	return weatherResponse, nil
}

func GetWeatherAPI(localidade string) (*models.WeatherResponse, error) {
	resp, err := http.Get("https://api.weatherapi.com/v1/current.json?key=" + os.Getenv("WEATHER_API_KEY") + "&q=" + localidade)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherResponse models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}

func ConvertTemperature(tempC float64) (float64, float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15
	return tempF, tempK, tempC
}

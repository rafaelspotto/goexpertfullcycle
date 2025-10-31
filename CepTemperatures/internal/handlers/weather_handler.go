package handlers

import (
	"ceptemperature/internal/models"
	"ceptemperature/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler para endpoint GET /weather/:cep
func GetWeatherByCep(c *gin.Context) {
	cep := c.Param("cep")

	// Validação do CEP
	isValid, err := services.ValidateCep(cep)
	if !isValid || err != nil {
		errorResp := models.ErrorResponse{
			Error: "invalid zipcode",
		}
		c.JSON(http.StatusUnprocessableEntity, errorResp)
		return
	}

	cepInfo, err := services.GetCep(cep)
	if err != nil {
		errorResp := models.ErrorResponse{
			Error: "invalid zipcode",
		}
		c.JSON(http.StatusUnprocessableEntity, errorResp)
		return
	}
	// Verificar se a ViaCEP retornou erro
	if cepInfo.Erro || cepInfo.Localidade == "" {
		errorResp := models.ErrorResponse{
			Error: "can not find zipcode",
		}
		c.JSON(http.StatusNotFound, errorResp)
		return
	}

	weather, err := services.GetWeather(cep)
	if err != nil {
		errorResp := models.ErrorResponse{
			Error: "can not find zipcode",
		}
		c.JSON(http.StatusNotFound, errorResp)
		return
	}

	// Retornar apenas os dados de temperatura conforme especificado
	c.JSON(http.StatusOK, weather)
}

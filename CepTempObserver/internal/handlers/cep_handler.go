package handlers

import (
	"ceptemperature/internal/models"
	"ceptemperature/internal/services"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var (
	tracer = otel.Tracer("service-a")
)

// HandleCepRequest handler para receber CEP via POST no Serviço A
func HandleCepRequest(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Criar span para o handler do Serviço A
	ctx, span := tracer.Start(ctx, "service-a.handle-cep-request")
	defer span.End()

	var req models.CepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid request format")
		
		errorResp := models.ErrorResponse{
			Error: "invalid zipcode",
		}
		c.JSON(http.StatusUnprocessableEntity, errorResp)
		return
	}

	span.SetAttributes(
		attribute.String("cep.input", req.Cep),
	)

	// Validar se o CEP é uma string válida com 8 dígitos
	isValid := validateCepString(req.Cep)
	if !isValid {
		err := fmt.Errorf("invalid zipcode")
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid CEP format")
		errorResp := models.ErrorResponse{
			Error: "invalid zipcode",
		}
		c.JSON(http.StatusUnprocessableEntity, errorResp)
		return
	}

	// Encaminhar para o Serviço B via HTTP com propagação de tracing
	response, err := services.ForwardToServiceB(ctx, req.Cep)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to forward request to Service B")
		
		// Determinar o código HTTP apropriado baseado no erro
		errMsg := err.Error()
		errorResp := models.ErrorResponse{
			Error: "can not find zipcode",
		}
		
		// Se for erro de validação, retornar 422
		if errMsg == "invalid zipcode" {
			errorResp.Error = "invalid zipcode"
			c.JSON(http.StatusUnprocessableEntity, errorResp)
			return
		}
		
		// Para outros erros, assumir 404 (CEP não encontrado)
		c.JSON(http.StatusNotFound, errorResp)
		return
	}

	span.SetStatus(codes.Ok, "Request processed successfully")
	c.JSON(http.StatusOK, response)
}

// validateCepString valida se o CEP é uma string válida com exatamente 8 dígitos
func validateCepString(cep string) bool {
	// Verificar se é string e contém exatamente 8 dígitos
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}



package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetWeatherByCep(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Valid CEP",
			cep:            "01310100",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid CEP - too short",
			cep:            "123",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "invalid zipcode",
		},
		{
			name:           "Invalid CEP - too long",
			cep:            "123456789",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "invalid zipcode",
		},
		{
			name:           "Invalid CEP - with letters",
			cep:            "1234567a",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "invalid zipcode",
		},
		{
			name:           "Empty CEP",
			cep:            "",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "invalid zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Gin router
			router := gin.New()
			router.GET("/weather/:cep", GetWeatherByCep)

			// Create a test request
			req, _ := http.NewRequest("GET", "/weather/"+tt.cep, nil)
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body for error cases
			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}

func TestGetWeatherByCepNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test with a CEP that doesn't exist
	router := gin.New()
	router.GET("/weather/:cep", GetWeatherByCep)

	req, _ := http.NewRequest("GET", "/weather/99999999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 404 for non-existent CEP
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", response["error"])
}

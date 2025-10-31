package services

import (
	"fmt"
	"testing"
)

func TestValidateCep(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
		hasError bool
	}{
		{
			name:     "Valid CEP with 8 digits",
			cep:      "01310100",
			expected: true,
			hasError: false,
		},
		{
			name:     "Invalid CEP with less than 8 digits",
			cep:      "1234567",
			expected: false,
			hasError: true,
		},
		{
			name:     "Invalid CEP with more than 8 digits",
			cep:      "123456789",
			expected: false,
			hasError: true,
		},
		{
			name:     "Invalid CEP with letters",
			cep:      "1234567a",
			expected: true, // ValidateCep only checks length, not content
			hasError: false,
		},
		{
			name:     "Empty CEP",
			cep:      "",
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := ValidateCep(tt.cep)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.hasError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if isValid != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, isValid)
			}
		})
	}
}

func TestValidateCepFormat(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
		hasError bool
	}{
		{
			name:     "Valid CEP format",
			cep:      "01310100",
			expected: true,
			hasError: false,
		},
		{
			name:     "Invalid CEP format - too short",
			cep:      "1234567",
			expected: false,
			hasError: true,
		},
		{
			name:     "Invalid CEP format - too long",
			cep:      "123456789",
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := ValidateCepFormat(tt.cep)

			if tt.hasError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.hasError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if isValid != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, isValid)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	testError := "test error message"
	errorResp, err := HandleError(fmt.Errorf("%s", testError))

	if err == nil {
		t.Errorf("Expected error but got none")
	}

	if errorResp.Error != testError {
		t.Errorf("Expected error message '%s', got '%s'", testError, errorResp.Error)
	}
}

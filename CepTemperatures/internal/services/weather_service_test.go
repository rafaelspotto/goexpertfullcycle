package services

import (
	"testing"
)

func TestConvertTemperature(t *testing.T) {
	tests := []struct {
		name      string
		tempC     float64
		expectedF float64
		expectedK float64
	}{
		{
			name:      "Convert 0°C",
			tempC:     0,
			expectedF: 32,
			expectedK: 273.15,
		},
		{
			name:      "Convert 25°C",
			tempC:     25,
			expectedF: 77,
			expectedK: 298.15,
		},
		{
			name:      "Convert 100°C",
			tempC:     100,
			expectedF: 212,
			expectedK: 373.15,
		},
		{
			name:      "Convert -40°C",
			tempC:     -40,
			expectedF: -40,
			expectedK: 233.15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempF, tempK, tempC := ConvertTemperature(tt.tempC)

			if tempC != tt.tempC {
				t.Errorf("Expected tempC %v, got %v", tt.tempC, tempC)
			}

			if tempF != tt.expectedF {
				t.Errorf("Expected tempF %v, got %v", tt.expectedF, tempF)
			}

			// Allow small floating point differences for Kelvin
			if abs(tempK-tt.expectedK) > 0.01 {
				t.Errorf("Expected tempK %v, got %v (difference: %v)", tt.expectedK, tempK, abs(tempK-tt.expectedK))
			}
		})
	}
}

func TestConvertTemperaturePrecision(t *testing.T) {
	// Test with decimal precision
	tempC := 28.5
	tempF, tempK, returnedTempC := ConvertTemperature(tempC)

	expectedF := 83.3
	expectedK := 301.65

	// Allow small floating point differences
	tolerance := 0.1

	if abs(tempF-expectedF) > tolerance {
		t.Errorf("Expected tempF %v, got %v (difference: %v)", expectedF, tempF, abs(tempF-expectedF))
	}

	if abs(tempK-expectedK) > tolerance {
		t.Errorf("Expected tempK %v, got %v (difference: %v)", expectedK, tempK, abs(tempK-expectedK))
	}

	if returnedTempC != tempC {
		t.Errorf("Expected tempC %v, got %v", tempC, returnedTempC)
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

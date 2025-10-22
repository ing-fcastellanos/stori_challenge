package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expected     int
		expectedFail bool
	}{
		{
			name:     "Valid integer",
			input:    "123",
			expected: 123,
		},
		{
			name:     "Negative integer",
			input:    "-123",
			expected: -123,
		},
		{
			name:     "Zero",
			input:    "0",
			expected: 0,
		},
		{
			name:         "Invalid input",
			input:        "abc",
			expected:     0,
			expectedFail: true, // Even though we get 0 here, we should assume that in real cases, this would fail gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseInt(tt.input)
			if tt.expectedFail {
				// Asserting that the result is not the expected value, because we want to handle errors
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expected     float64
		expectedFail bool
	}{
		{
			name:     "Valid float",
			input:    "123.45",
			expected: 123.45,
		},
		{
			name:     "Negative float",
			input:    "-123.45",
			expected: -123.45,
		},
		{
			name:     "Zero float",
			input:    "0.0",
			expected: 0.0,
		},
		{
			name:         "Invalid float",
			input:        "abc",
			expected:     0.0,
			expectedFail: true, // Similar to the integer test, we want to handle invalid cases as gracefully as possible
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFloat(tt.input)
			if tt.expectedFail {
				// Asserting that the result is not the expected value
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

package ibm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanRegion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "region with zone suffix eu-de-1",
			input:    "eu-de-1",
			expected: "eu-de",
		},
		{
			name:     "region with zone suffix eu-de-2",
			input:    "eu-de-2",
			expected: "eu-de",
		},
		{
			name:     "region with zone suffix us-south-1",
			input:    "us-south-1",
			expected: "us-south",
		},
		{
			name:     "region with custom suffix eu-de-zone1",
			input:    "eu-de-zone1",
			expected: "eu-de",
		},
		{
			name:     "standard region without suffix eu-de",
			input:    "eu-de",
			expected: "eu-de",
		},
		{
			name:     "standard region without suffix us-south",
			input:    "us-south",
			expected: "us-south",
		},
		{
			name:     "single token region global",
			input:    "global",
			expected: "global",
		},
		{
			name:     "empty region",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CleanRegion(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

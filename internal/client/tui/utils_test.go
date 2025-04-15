package tui_test

import (
	"testing"

	"github.com/npavlov/go-password-manager/internal/client/tui"
	"github.com/stretchr/testify/require"
)

func TestFormatCardNumber(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic 16-digit",
			input:    "1234567812345678",
			expected: "1234 5678 1234 5678",
		},
		{
			name:     "16-digit with spaces",
			input:    "1234 5678 1234 5678",
			expected: "1234 5678 1234 5678",
		},
		{
			name:     "with dashes and spaces",
			input:    "1234-5678-1234-5678",
			expected: "1234 5678 1234 5678",
		},
		{
			name:     "non-digit characters",
			input:    "abcd1234!@#5678efgh9012ijkl3456",
			expected: "1234 5678 9012 3456",
		},
		{
			name:     "short input",
			input:    "123",
			expected: "123",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := tui.FormatCardNumber(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

package utils_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

func TestHashCardNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "valid card number",
			input: "4111111111111111",
		},
		{
			name:  "empty card number",
			input: "",
		},
		{
			name:  "card with dashes",
			input: "1234-5678-9012-3456",
		},
		{
			name:  "long card number",
			input: "55554444333322221111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expected SHA-256 hash as hex string
			hash := sha256.Sum256([]byte(tt.input))
			expectedHash := hex.EncodeToString(hash[:])

			// Call the function
			result := utils.HashCardNumber(tt.input)

			// Check that result is a valid pgtype.Text
			if !result.Valid {
				t.Errorf("expected Valid=true but got false")
			}

			// Check the actual hash value
			if result.String != expectedHash {
				t.Errorf("expected hash %s but got %s", expectedHash, result.String)
			}
		})
	}
}

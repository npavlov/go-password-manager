package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgtype"
)

// HashCardNumber hashes the card number to enforce uniqueness.
func HashCardNumber(cardNumber string) pgtype.Text {
	hash := sha256.Sum256([]byte(cardNumber))

	text := pgtype.Text{
		String: "",
		Valid:  false,
	}

	hashString := hex.EncodeToString(hash[:]) // Convert to hex string for storage

	_ = text.Scan(hashString)

	return text
}

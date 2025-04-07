package tui

import "strings"

func FormatCardNumber(cardNumber string) string {
	// Remove all non-digit characters
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, cardNumber)

	// Split into groups of 4 digits
	var parts []string
	for i := 0; i < len(cleaned); i += 4 {
		end := i + 4
		if end > len(cleaned) {
			end = len(cleaned)
		}
		parts = append(parts, cleaned[i:end])
	}

	return strings.Join(parts, " ")
}

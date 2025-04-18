package model

import (
	"time"
)

// ItemType defines the type of stored data.
type ItemType string

const (
	ItemTypePassword ItemType = "password"
	ItemTypeNote     ItemType = "note"
	ItemTypeCard     ItemType = "card"
	ItemTypeBinary   ItemType = "binary"
)

// StorageItem represents a generic stored item with metadata.
type StorageItem struct {
	ID        string            `json:"id"`
	Type      ItemType          `json:"type"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Metadata  map[string]string `json:"metadata"` // Key-value metadata
}

// PasswordItem stores encrypted passwords.
type PasswordItem struct {
	StorageItem
	Login    string `json:"username"`
	Password string `json:"password"` // Encrypted
}

// NoteItem stores secure notes.
type NoteItem struct {
	StorageItem
	Content string `json:"content"` // Encrypted
}

// CardItem stores encrypted card details.
type CardItem struct {
	StorageItem
	CardNumber     string `json:"cardNumber"` // Encrypted
	ExpiryDate     string `json:"expiryDate"`
	CVV            string `json:"cvv"` // Encrypted
	CardholderName string `json:"cardholderName"`
}

// BinaryItem stores metadata for large binary files.
type BinaryItem struct {
	StorageItem
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
}

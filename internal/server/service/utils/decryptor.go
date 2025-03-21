package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Decryptor wraps an io.Reader and decrypts data in blocks
type Decryptor struct {
	reader io.Reader
	gcm    cipher.AEAD
}

func NewDecryptor(r io.Reader, base64Key string) (*Decryptor, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Decryptor{reader: r, gcm: gcm}, nil
}

func (d *Decryptor) Read(p []byte) (int, error) {
	nonceSize := d.gcm.NonceSize()
	nonce := make([]byte, nonceSize)

	// Read the nonce first
	if _, err := io.ReadFull(d.reader, nonce); err != nil {
		return 0, fmt.Errorf("failed to read nonce: %v", err)
	}

	// Read the encrypted block
	encrypted := make([]byte, len(p)+d.gcm.Overhead())
	n, err := d.reader.Read(encrypted)
	if err != nil && err != io.EOF {
		return 0, err
	}

	// Decrypt the data
	decrypted, err := d.gcm.Open(nil, nonce, encrypted[:n], nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to decrypt data")
	}

	copy(p, decrypted)

	return len(decrypted), err
}

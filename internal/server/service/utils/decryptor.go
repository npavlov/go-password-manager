//nolint:wrapcheck
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
)

// Decryptor wraps an io.Reader and decrypts data in blocks.
type Decryptor struct {
	reader io.Reader
	gcm    cipher.AEAD
}

func NewDecryptor(reader io.Reader, base64Key string) (*Decryptor, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Decryptor{reader: reader, gcm: gcm}, nil
}

func (d *Decryptor) Read(bytes []byte) (int, error) {
	nonceSize := d.gcm.NonceSize()
	nonce := make([]byte, nonceSize)

	// Read the nonce first
	if _, err := io.ReadFull(d.reader, nonce); err != nil {
		return 0, err
	}

	// Read the encrypted block
	encrypted := make([]byte, len(bytes)+d.gcm.Overhead())
	cursor, err := d.reader.Read(encrypted)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}

	// Decrypt the data
	decrypted, err := d.gcm.Open(nil, nonce, encrypted[:cursor], nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to decrypt data")
	}

	copy(bytes, decrypted)

	return len(decrypted), err
}

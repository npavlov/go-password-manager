//nolint:wrapcheck
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Encryptor wraps an io.Writer and encrypts data in blocks.
type Encryptor struct {
	writer io.Writer
	block  cipher.Block
	gcm    cipher.AEAD
}

func NewEncryptor(writer io.Writer, base64Key string) (*Encryptor, error) {
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

	return &Encryptor{writer: writer, block: block, gcm: gcm}, nil
}

func (e *Encryptor) Write(bytes []byte) (int, error) {
	blockSize := 1024
	offset := 0
	totalLength := 0

	for offset < len(bytes) {
		end := offset + blockSize
		if end > len(bytes) {
			end = len(bytes)
		}

		nonce := make([]byte, e.gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return offset, err
		}

		encrypted := e.gcm.Seal(nil, nonce, bytes[offset:end], nil)

		// Write nonce + encrypted block
		if _, err := e.writer.Write(nonce); err != nil {
			return offset, err
		}
		n, err := e.writer.Write(encrypted)
		if err != nil {
			return offset + n, err
		}
		totalLength += len(nonce) + n

		offset += blockSize
	}

	return totalLength, nil
}

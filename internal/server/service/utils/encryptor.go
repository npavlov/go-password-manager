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

func NewEncryptor(w io.Writer, base64Key string) (*Encryptor, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Encryptor{writer: w, block: block, gcm: gcm}, nil
}

func (e *Encryptor) Write(p []byte) (int, error) {
	blockSize := 1024
	offset := 0
	totalLength := 0

	for offset < len(p) {
		end := offset + blockSize
		if end > len(p) {
			end = len(p)
		}

		nonce := make([]byte, e.gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return offset, err
		}

		encrypted := e.gcm.Seal(nil, nonce, p[offset:end], nil)

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

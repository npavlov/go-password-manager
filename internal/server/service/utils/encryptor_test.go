package utils_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

func generateTestKey(t *testing.T) string {
	t.Helper()

	key := make([]byte, 32) // AES-256
	_, err := rand.Read(key)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(key)
}

func TestEncryptor_Write_Success(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	key := generateTestKey(t)

	encryptor, err := utils.NewEncryptor(&buf, key)
	require.NoError(t, err, "NewEncryptor should not return an error")

	plaintext := []byte("This is a test of the emergency encryption system.")

	// Act
	n, err := encryptor.Write(plaintext)

	// Assert
	require.NoError(t, err, "Write should not return an error")
	assert.Positive(t, n, "Write should write more than 0 bytes")
	assert.Positive(t, buf.Len(), "Buffer should contain encrypted data")
	assert.NotEqual(t, plaintext, buf.Bytes(), "Encrypted data should not match plaintext")
}

func TestEncryptor_Write_LongInput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	key := generateTestKey(t)

	encryptor, err := utils.NewEncryptor(&buf, key)
	require.NoError(t, err)

	// 5000 bytes to span multiple 1024-byte blocks
	plaintext := make([]byte, 5000)
	_, err = rand.Read(plaintext)
	require.NoError(t, err)

	// Act
	n, err := encryptor.Write(plaintext)

	// Assert
	require.NoError(t, err)
	assert.Positive(t, n)
	assert.Positive(t, buf.Len())
}

func TestEncryptor_InvalidKey(t *testing.T) {
	// Arrange
	var buf bytes.Buffer

	// Invalid base64 (not even decodable)
	_, err := utils.NewEncryptor(&buf, "invalid-base64!")
	require.Error(t, err, "NewEncryptor should return error for invalid base64 key")

	// Valid base64, but wrong length (e.g., 10 bytes instead of 16, 24, or 32)
	badKey := base64.StdEncoding.EncodeToString([]byte("short-key"))
	_, err = utils.NewEncryptor(&buf, badKey)
	require.Error(t, err, "NewEncryptor should return error for short key")
}

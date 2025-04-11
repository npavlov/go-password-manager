package utils_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

func TestGenerateRandomKey(t *testing.T) {
	key1, err1 := utils.GenerateRandomKey()
	key2, err2 := utils.GenerateRandomKey()

	assert.NoError(t, err1, "should not error generating key1")
	assert.NoError(t, err2, "should not error generating key2")
	assert.NotEmpty(t, key1, "key1 should not be empty")
	assert.NotEmpty(t, key2, "key2 should not be empty")
	assert.NotEqual(t, key1, key2, "keys should be different")

	// Check that key decodes to 32 bytes
	decodedKey, err := base64.StdEncoding.DecodeString(key1)
	assert.NoError(t, err)
	assert.Len(t, decodedKey, 32, "AES-256 key should be 32 bytes")
}

func TestEncryptDecrypt_Success(t *testing.T) {
	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	originalText := "Secret message for encryption"

	encrypted, err := utils.Encrypt(originalText, key)
	require.NoError(t, err)
	assert.NotEmpty(t, encrypted, "encrypted text should not be empty")

	decrypted, err := utils.Decrypt(encrypted, key)
	require.NoError(t, err)
	assert.Equal(t, originalText, decrypted, "decrypted text should match original")
}

func TestEncrypt_InvalidKey(t *testing.T) {
	invalidKey := "short-key"
	_, err := utils.Encrypt("data", invalidKey)
	assert.Error(t, err, "encryption should fail with invalid key")
}

func TestDecrypt_InvalidKey(t *testing.T) {
	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	encrypted, err := utils.Encrypt("test", key)
	require.NoError(t, err)

	// Use a different invalid key for decryption
	invalidKey := "invalid-base64=="
	_, err = utils.Decrypt(encrypted, invalidKey)
	assert.Error(t, err, "decryption should fail with invalid base64 key")
}

func TestDecrypt_TamperedData(t *testing.T) {
	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	encrypted, err := utils.Encrypt("hello", key)
	require.NoError(t, err)

	// Tamper with the ciphertext (e.g., flip a byte)
	tampered := []byte(encrypted)
	tampered[len(tampered)-1] ^= 0xFF

	_, err = utils.Decrypt(string(tampered), key)
	assert.Error(t, err, "decryption should fail on tampered data")
}

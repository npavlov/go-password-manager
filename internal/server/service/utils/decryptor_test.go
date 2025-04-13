package utils_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

func TestDecryptor_Read_Success(t *testing.T) {
	// Arrange
	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	originalText := []byte("streaming test message across multiple blocks")
	var encryptedBuf bytes.Buffer

	encryptor, err := utils.NewEncryptor(&encryptedBuf, key)
	require.NoError(t, err)

	n, err := encryptor.Write(originalText)
	require.NoError(t, err)
	require.Positive(t, n)

	// Act
	decryptor, err := utils.NewDecryptor(&encryptedBuf, key)
	require.NoError(t, err)

	decrypted := make([]byte, len(originalText))
	read, err := decryptor.Read(decrypted)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, originalText, decrypted[:read])
}

func TestDecryptor_Read_InvalidKey(t *testing.T) {
	// Arrange: valid encrypted data with valid key
	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	var buf bytes.Buffer
	encryptor, err := utils.NewEncryptor(&buf, key)
	require.NoError(t, err)

	_, err = encryptor.Write([]byte("hello"))
	require.NoError(t, err)

	// Tamper with the key (wrong key)
	invalidKey := base64.StdEncoding.EncodeToString([]byte("wrong-secret-wrong-secret-wrong-secret!!"))

	// Act
	_, err = utils.NewDecryptor(&buf, invalidKey)
	require.Error(t, err)
}

func TestDecryptor_Read_TamperedData(t *testing.T) {
	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	var buf bytes.Buffer
	encryptor, err := utils.NewEncryptor(&buf, key)
	require.NoError(t, err)

	_, err = encryptor.Write([]byte("important data"))
	require.NoError(t, err)

	encrypted := buf.Bytes()
	// Flip a byte in ciphertext to simulate tampering
	encrypted[len(encrypted)-1] ^= 0xFF

	tamperedBuf := bytes.NewBuffer(encrypted)

	decryptor, err := utils.NewDecryptor(tamperedBuf, key)
	require.NoError(t, err)

	out := make([]byte, 128)
	_, err = decryptor.Read(out)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt data")
}

func TestDecryptor_Read_InvalidNonce(t *testing.T) {
	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	// Create a buffer with insufficient nonce size
	badData := bytes.NewBufferString("shortnonce")

	decryptor, err := utils.NewDecryptor(badData, key)
	require.NoError(t, err)

	buf := make([]byte, 64)
	_, err = decryptor.Read(buf)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected EOF")
}

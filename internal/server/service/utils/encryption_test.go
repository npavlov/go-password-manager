package utils_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/service/utils"
)

func TestGenerateRandomKey(t *testing.T) {
	t.Parallel()

	t.Run("should generate valid base64 encoded key", func(t *testing.T) {
		t.Parallel()

		key, err := utils.GenerateRandomKey()
		require.NoError(t, err)
		assert.NotEmpty(t, key)

		// Verify it's valid base64
		decoded, err := base64.StdEncoding.DecodeString(key)
		require.NoError(t, err)
		assert.Len(t, decoded, 32) // AES-256 key size
	})

	t.Run("should generate unique keys each time", func(t *testing.T) {
		t.Parallel()

		key1, err := utils.GenerateRandomKey()
		require.NoError(t, err)

		key2, err := utils.GenerateRandomKey()
		require.NoError(t, err)

		assert.NotEqual(t, key1, key2)
	})
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		text string
	}{
		{"empty string", ""},
		{"short text", "hello world"},
		{"long text", "Lorem ipsum dolor sit amet, consectetur adipiscing elit."},
		{"special characters", "!@#$%^&*()_+-=[]{};':\",./<>?"},
		{"unicode text", "こんにちは世界"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			key, err := utils.GenerateRandomKey()
			require.NoError(t, err)

			// Test encryption
			encrypted, err := utils.Encrypt(tc.text, key)
			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)
			assert.NotEqual(t, tc.text, encrypted)

			// Test decryption
			decrypted, err := utils.Decrypt(encrypted, key)
			require.NoError(t, err)
			assert.Equal(t, tc.text, decrypted)
		})
	}

	t.Run("should fail with invalid key", func(t *testing.T) {
		t.Parallel()

		_, err := utils.Encrypt("test", "not a valid base64 key")
		require.Error(t, err)

		_, err = utils.Decrypt("encrypted text", "not a valid base64 key")
		require.Error(t, err)
	})

	t.Run("should fail with wrong key", func(t *testing.T) {
		t.Parallel()

		key1, err := utils.GenerateRandomKey()
		require.NoError(t, err)

		key2, err := utils.GenerateRandomKey()
		require.NoError(t, err)

		encrypted, err := utils.Encrypt("test message", key1)
		require.NoError(t, err)

		_, err = utils.Decrypt(encrypted, key2)
		require.Error(t, err)
	})

	t.Run("should fail with tampered ciphertext", func(t *testing.T) {
		t.Parallel()

		key, err := utils.GenerateRandomKey()
		require.NoError(t, err)

		encrypted, err := utils.Encrypt("test message", key)
		require.NoError(t, err)

		// Tamper with the ciphertext
		decoded, err := base64.StdEncoding.DecodeString(encrypted)
		require.NoError(t, err)

		// Change one byte in the ciphertext
		if len(decoded) > 0 {
			decoded[0] ^= 0xFF
		}

		tampered := base64.StdEncoding.EncodeToString(decoded)

		_, err = utils.Decrypt(tampered, key)
		require.Error(t, err)
	})

	t.Run("should fail with short ciphertext", func(t *testing.T) {
		t.Parallel()

		key, err := utils.GenerateRandomKey()
		require.NoError(t, err)

		_, err = utils.Decrypt("a", key) // way too short
		require.Error(t, err)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty key", func(t *testing.T) {
		t.Parallel()

		_, err := utils.Encrypt("test", "")
		require.Error(t, err)

		_, err = utils.Decrypt("test", "")
		require.Error(t, err)
	})

	t.Run("nil inputs", func(t *testing.T) {
		t.Parallel()

		key, _ := utils.GenerateRandomKey()
		_, err := utils.Encrypt("", key) // empty text is okay
		require.NoError(t, err)

		_, err = utils.Encrypt("test", "") // empty key is not
		require.Error(t, err)
	})

	t.Run("decrypt non-base64 text", func(t *testing.T) {
		t.Parallel()

		key, err := utils.GenerateRandomKey()
		require.NoError(t, err)

		_, err = utils.Decrypt("not base64 encoded", key)
		require.Error(t, err)
	})

	t.Run("decrypt with wrong key size", func(t *testing.T) {
		// Create a key with wrong size (not 32 bytes)
		wrongKey := base64.StdEncoding.EncodeToString([]byte("short key"))

		_, err := utils.Encrypt("test", wrongKey)
		require.Error(t, err)

		_, err = utils.Decrypt("some encrypted text", wrongKey)
		require.Error(t, err)
	})
}

func TestEncryptDecryptConsistency(t *testing.T) {
	t.Parallel()

	key, err := utils.GenerateRandomKey()
	require.NoError(t, err)

	originalText := "This text should remain the same after encryption and decryption"

	// Run multiple times to ensure consistency
	for range 10 {
		encrypted, err := utils.Encrypt(originalText, key)
		require.NoError(t, err)

		decrypted, err := utils.Decrypt(encrypted, key)
		require.NoError(t, err)

		assert.Equal(t, originalText, decrypted)
	}
}

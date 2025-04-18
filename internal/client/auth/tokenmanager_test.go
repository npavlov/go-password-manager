//nolint:exhaustruct,ireturn
package auth_test

import (
	"encoding/json"
	"os"
	"testing"

	obs "github.com/Dentrax/obscure-go/observer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/client/auth"
	"github.com/npavlov/go-password-manager/internal/client/config"
	"github.com/npavlov/go-password-manager/internal/utils"
)

type MockSecureString struct {
	value string
}

func (m *MockSecureString) Apply() utils.ISecureString {
	panic("implement me")
}

func (m *MockSecureString) AddWatcher(_ obs.Observer) {
	panic("implement me")
}

func (m *MockSecureString) SetKey(_ int) {
	panic("implement me")
}

func (m *MockSecureString) GetSelf() *utils.SecureString {
	panic("implement me")
}

func (m *MockSecureString) Decrypt() []rune {
	panic("implement me")
}

func (m *MockSecureString) RandomizeKey() {
	panic("implement me")
}

func (m *MockSecureString) IsEquals(_ utils.ISecureString) bool {
	panic("implement me")
}

func (m *MockSecureString) Get() string {
	return m.value
}

func (m *MockSecureString) Set(value string) utils.ISecureString {
	m.value = value

	return m
}

func (m *MockSecureString) Clear() {
	m.value = ""
}

func TestNewTokenManager(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "test_tokens.json"}

	tm := auth.NewTokenManager(&logger, cfg)

	tm.AccessToken = &MockSecureString{
		value: "access_token",
	}
	tm.RefreshToken = &MockSecureString{
		value: "refresh_token",
	}

	assert.NotNil(t, tm)
	assert.False(t, tm.IsAuthorized())
}

func TestLoadTokens_Success(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "test_tokens.json"}

	// Create test token file
	testTokens := auth.DataObject{
		AccessToken:  "test_access",
		RefreshToken: "test_refresh",
	}
	data, err := json.Marshal(testTokens)
	require.NoError(t, err)
	err = os.WriteFile(cfg.TokenFile, data, 0o600)
	require.NoError(t, err)
	defer os.Remove(cfg.TokenFile)

	tm := auth.NewTokenManager(&logger, cfg)
	err = tm.LoadTokens()
	require.NoError(t, err)
	assert.True(t, tm.IsAuthorized())
	assert.Equal(t, "test_access", tm.GetAccessToken())
	assert.Equal(t, "test_refresh", tm.GetRefreshToken())
}

func TestLoadTokens_FileNotExist(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "nonexistent.json"}

	tm := auth.NewTokenManager(&logger, cfg)
	err := tm.LoadTokens()
	require.NoError(t, err)
	assert.False(t, tm.IsAuthorized())
}

func TestLoadTokens_InvalidJSON(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "invalid_tokens.json"}

	// Create invalid JSON file
	err := os.WriteFile(cfg.TokenFile, []byte("{invalid json}"), 0o600)
	require.NoError(t, err)
	defer os.Remove(cfg.TokenFile)

	tm := auth.NewTokenManager(&logger, cfg)
	err = tm.LoadTokens()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode token")
}

func TestSaveTokens_Success(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "save_test.json"}
	defer os.Remove(cfg.TokenFile)

	tm := auth.NewTokenManager(&logger, cfg)
	err := tm.SaveTokens("new_access", "new_refresh")
	require.NoError(t, err)

	// Verify file contents
	data, err := os.ReadFile(cfg.TokenFile)
	require.NoError(t, err)

	var tokens auth.DataObject
	err = json.Unmarshal(data, &tokens)
	require.NoError(t, err)
	assert.Equal(t, "new_access", tokens.AccessToken)
	assert.Equal(t, "new_refresh", tokens.RefreshToken)
}

func TestSaveTokens_FileError(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "/invalid/path/tokens.json"}

	tm := auth.NewTokenManager(&logger, cfg)
	err := tm.SaveTokens("access", "refresh")
	require.Error(t, err)
}

func TestUpdateTokens(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "update_test.json"}
	defer os.Remove(cfg.TokenFile)

	tm := auth.NewTokenManager(&logger, cfg)
	err := tm.UpdateTokens("updated_access", "updated_refresh")
	require.NoError(t, err)

	assert.Equal(t, "updated_access", tm.GetAccessToken())
	assert.Equal(t, "updated_refresh", tm.GetRefreshToken())

	// Verify tokens were saved to file
	data, err := os.ReadFile(cfg.TokenFile)
	require.NoError(t, err)

	var tokens auth.DataObject
	err = json.Unmarshal(data, &tokens)
	require.NoError(t, err)
	assert.Equal(t, "updated_access", tokens.AccessToken)
	assert.Equal(t, "updated_refresh", tokens.RefreshToken)
}

func TestHandleAuthFailure(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "auth_failure_test.json"}
	defer os.Remove(cfg.TokenFile)

	// Setup initial tokens
	tm := auth.NewTokenManager(&logger, cfg)
	err := tm.UpdateTokens("valid_access", "valid_refresh")
	require.NoError(t, err)

	// Set up auth failure callback
	called := false
	tm.SetAuthFailCallback(func() {
		called = true
		require.NoError(t, err)
	})

	tm.HandleAuthFailure()

	assert.False(t, tm.IsAuthorized())
	assert.True(t, called)

	// Verify tokens were cleared from file
	data, err := os.ReadFile(cfg.TokenFile)
	require.NoError(t, err)

	var tokens auth.DataObject
	err = json.Unmarshal(data, &tokens)
	require.NoError(t, err)
	assert.Empty(t, tokens.AccessToken)
	assert.Empty(t, tokens.RefreshToken)
}

func TestGetTokens(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()
	cfg := &config.Config{TokenFile: "get_tokens_test.json"}

	tm := auth.NewTokenManager(&logger, cfg)
	tm.AccessToken = &MockSecureString{value: "test_access"}
	tm.RefreshToken = &MockSecureString{value: "test_refresh"}

	assert.Equal(t, "test_access", tm.GetAccessToken())
	assert.Equal(t, "test_refresh", tm.GetRefreshToken())
}

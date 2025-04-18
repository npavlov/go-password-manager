package auth

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/npavlov/go-password-manager/internal/client/config"
	"github.com/npavlov/go-password-manager/internal/utils"
)

type TokenManager struct {
	mu           sync.Mutex
	logger       *zerolog.Logger
	isAuthorized bool
	onAuthFail   func() // UI callback function
	AccessToken  utils.ISecureString
	RefreshToken utils.ISecureString
	cfg          *config.Config
}

type DataObject struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewTokenManager(logger *zerolog.Logger, cfg *config.Config) *TokenManager {
	return &TokenManager{
		logger:       logger,
		isAuthorized: false,
		mu:           sync.Mutex{},
		onAuthFail:   func() {},
		AccessToken:  nil,
		RefreshToken: nil,
		cfg:          cfg,
	}
}

// LoadTokens reads the tokens from file.
func (tm *TokenManager) LoadTokens() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	file, err := os.Open(tm.cfg.TokenFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // No file means first-time login
		}

		return errors.Wrap(err, "failed to open token file")
	}
	defer file.Close()

	dataObj := DataObject{
		AccessToken:  "",
		RefreshToken: "",
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dataObj); err != nil {
		return errors.Wrap(err, "failed to decode token")
	}

	tm.isAuthorized = dataObj.AccessToken != "" && dataObj.RefreshToken != ""
	tm.AccessToken = utils.NewString(dataObj.AccessToken)
	tm.RefreshToken = utils.NewString(dataObj.RefreshToken)
	tm.isAuthorized = true

	return nil
}

// SaveTokens writes tokens to a file.
func (tm *TokenManager) SaveTokens(accessToken, refreshToken string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	data, err := json.MarshalIndent(&DataObject{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to serialize token")
	}

	//nolint:mnd,wrapcheck
	return os.WriteFile(tm.cfg.TokenFile, data, 0o600)
}

// UpdateTokens saves new tokens and marks user as authorized.
func (tm *TokenManager) UpdateTokens(access, refresh string) error {
	tm.AccessToken = utils.NewString(access)
	tm.RefreshToken = utils.NewString(refresh)
	tm.isAuthorized = true

	return tm.SaveTokens(access, refresh)
}

// IsAuthorized returns the user's authentication state.
func (tm *TokenManager) IsAuthorized() bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	return tm.isAuthorized
}

// HandleAuthFailure clears tokens and notifies the UI.
func (tm *TokenManager) HandleAuthFailure() {
	tm.AccessToken = nil
	tm.RefreshToken = nil
	tm.isAuthorized = false
	err := tm.SaveTokens("", "")
	if err != nil {
		return
	}

	if tm.onAuthFail != nil {
		tm.onAuthFail() // Notify UI to go back to login screen
	}
}

// SetAuthFailCallback sets a UI callback when authentication fails.
func (tm *TokenManager) SetAuthFailCallback(callback func()) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.onAuthFail = callback
}

func (tm *TokenManager) GetAccessToken() string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	return tm.AccessToken.Get()
}

func (tm *TokenManager) GetRefreshToken() string {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	return tm.RefreshToken.Get()
}

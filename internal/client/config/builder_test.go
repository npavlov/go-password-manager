//nolint:exhaustruct,paralleltest
package config_test

import (
	"flag"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/npavlov/go-password-manager/internal/client/config"
)

func TestNewConfigBuilder(t *testing.T) {
	logger := zerolog.Nop()
	builder := config.NewConfigBuilder(&logger)

	assert.NotNil(t, builder)
}

func TestFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *config.Config
	}{
		{
			name: "All environment variables set",
			envVars: map[string]string{
				"ADDRESS":     "localhost:8080",
				"MASTER_KEY":  "testkey",
				"CERTIFICATE": "cert.pem",
				"TOKEN_FILE":  "tokens.json",
			},
			expected: &config.Config{
				Address:     "localhost:8080",
				MasterKey:   "testkey",
				Certificate: "cert.pem",
				TokenFile:   "tokens.json",
			},
		},
		{
			name:    "No environment variables set",
			envVars: map[string]string{},
			expected: &config.Config{
				Address:     ":9090", // default from envDefault tag
				MasterKey:   "",
				Certificate: "",
				TokenFile:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			logger := zerolog.Nop()
			builder := config.NewConfigBuilder(&logger).FromEnv()
			cfg := builder.Build()

			assert.Equal(t, tt.expected.Address, cfg.Address)
			assert.Equal(t, tt.expected.Certificate, cfg.Certificate)
			assert.Equal(t, tt.expected.TokenFile, cfg.TokenFile)
			// MasterKey should be cleared after Build()
			assert.Empty(t, cfg.MasterKey)
			assert.NotNil(t, cfg.SecuredMasterKey)
		})
	}
}

func TestFromFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *config.Config
	}{
		{
			name: "All flags set",
			args: []string{
				"-a", "localhost:8080",
				"-masterkey", "testkey",
				"-cert", "cert.pem",
				"-token_file", "tokens.json",
			},
			expected: &config.Config{
				Address:     "localhost:8080",
				MasterKey:   "testkey",
				Certificate: "cert.pem",
				TokenFile:   "tokens.json",
			},
		},
		{
			name: "No flags set",
			args: []string{},
			expected: &config.Config{
				Address:     "",
				MasterKey:   "",
				Certificate: "",
				TokenFile:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine to avoid flag redefinition errors
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)

			logger := zerolog.Nop()
			builder := config.NewConfigBuilder(&logger)

			// Simulate command line arguments
			os.Args = append([]string{"test"}, tt.args...)
			builder.FromFlags()
			cfg := builder.Build()

			assert.Equal(t, tt.expected.Address, cfg.Address)
			assert.Equal(t, tt.expected.Certificate, cfg.Certificate)
			assert.Equal(t, tt.expected.TokenFile, cfg.TokenFile)
			// MasterKey should be cleared after Build()
			assert.Empty(t, cfg.MasterKey)
			assert.NotNil(t, cfg.SecuredMasterKey)
		})
	}
}

func TestFromObj(t *testing.T) {
	t.Parallel()

	inputCfg := &config.Config{
		Address:     "localhost:8080",
		MasterKey:   "testkey",
		Certificate: "cert.pem",
		TokenFile:   "tokens.json",
	}

	logger := zerolog.Nop()
	builder := config.NewConfigBuilder(&logger).FromObj(inputCfg)
	cfg := builder.Build()

	assert.Equal(t, inputCfg.Address, cfg.Address)
	assert.Equal(t, inputCfg.Certificate, cfg.Certificate)
	assert.Equal(t, inputCfg.TokenFile, cfg.TokenFile)
	// MasterKey should be cleared after Build()
	assert.Empty(t, cfg.MasterKey)
	assert.NotNil(t, cfg.SecuredMasterKey)
}

func TestBuild(t *testing.T) {
	tests := []struct {
		name          string
		masterKey     string
		expectSecure  bool
		expectCleared bool
	}{
		{
			name:          "With master key",
			masterKey:     "testkey",
			expectSecure:  true,
			expectCleared: true,
		},
		{
			name:          "Empty master key",
			masterKey:     "",
			expectSecure:  false,
			expectCleared: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			builder := config.NewConfigBuilder(&logger)

			if tt.masterKey != "" {
				t.Setenv("MASTER_KEY", tt.masterKey)
			}

			cfg := builder.Build()

			assert.Empty(t, cfg.MasterKey)

			assert.Empty(t, cfg.MasterKey)
		})
	}
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/client/config"
	"github.com/npavlov/go-password-manager/internal/client/interceptors"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestLoadConfig(t *testing.T) {
	t.Run("successful config load", func(t *testing.T) {
		// Setup
		log := testutils.GetTLogger()

		// Set test environment variables
		t.Setenv("ADDRESS", "localhost:50051")
		t.Setenv("CERTIFICATE", "cert.pem")

		// Execute
		cfg := LoadConfig(log)

		// Verify
		assert.NotNil(t, cfg)
		assert.Equal(t, "localhost:50051", cfg.Address)
		assert.Equal(t, "cert.pem", cfg.Certificate)
	})
}

func TestMakeConnection(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		// Setup
		cfg := config.Config{
			Address:     "localhost:50051",
			Certificate: "testdata/cert.pem", // You'll need a test certificate file
		}

		interceptor := &interceptors.AuthInterceptor{}

		// Execute
		conn, err := MakeConnection(cfg, interceptor)

		// Verify
		require.NoError(t, err)
		assert.NotNil(t, conn)

		if conn != nil {
			conn.Close()
		}
	})

	t.Run("invalid certificate", func(t *testing.T) {
		// Setup
		cfg := config.Config{
			Address:     "localhost:50051",
			Certificate: "invalid-cert.pem",
		}

		interceptor := &interceptors.AuthInterceptor{}

		// Execute
		conn, err := MakeConnection(cfg, interceptor)

		// Verify
		require.Error(t, err)
		assert.Nil(t, conn)
		assert.Contains(t, err.Error(), "could not load TLS keys")
	})
}

// Note: Testing GetApp and GetTUI would require extensive mocking of all dependencies.
// In practice, you might want to test these components separately through their own packages.

func TestGetApp(t *testing.T) {
	logger := testutils.GetTLogger()

	t.Setenv("ADDRESS", "localhost:50051")
	t.Setenv("CERTIFICATE", "testdata/cert.pem")

	tokenMgr, facade, stMgr, gRPC := GetApp(logger)
	assert.NotNil(t, tokenMgr)
	assert.NotNil(t, facade)
	assert.NotNil(t, stMgr)
	assert.NotNil(t, gRPC)
}

func TestGetTUI(t *testing.T) {
	logger := testutils.GetTLogger()

	t.Setenv("ADDRESS", "localhost:50051")
	t.Setenv("CERTIFICATE", "testdata/cert.pem")

	tokenMgr, facade, stMgr, _ := GetApp(logger)

	tui := GetTUI(logger, facade, stMgr, tokenMgr)

	require.NotNil(t, tui)
}

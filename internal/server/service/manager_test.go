package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/service"
	testutils "github.com/npavlov/go-password-manager/internal/test_utils"
)

func TestGManager_StartAndShutdown(t *testing.T) {
	t.Parallel()

	// Arrange
	logger := testutils.GetTLogger()

	cfg := &config.Config{
		Address:     "localhost:50051",
		Certificate: "testdata/cert.pem", // use valid test certs or mock creds
		PrivateKey:  "testdata/key.pem",
		JwtSecret:   "testsecret",
	}

	mockRedis := testutils.NewMockRedis()

	// Use a wait group to wait for shutdown
	wg := &sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(t.Context())

	gm := service.NewGRPCManager(cfg, logger, mockRedis)

	// Act
	go gm.Start(ctx, wg)

	// Give some time for the server to start (could use readiness probe in real tests)
	time.Sleep(500 * time.Millisecond)

	// Simulate shutdown
	cancel()

	// Wait for graceful shutdown
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}

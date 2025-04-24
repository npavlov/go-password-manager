//nolint:err113,goconst,exhaustruct
package redis_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redismock/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/redis"
)

func TestNewRStorage(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(nil)
	cfg := config.Config{
		Redis: "localhost:6379",
	}

	storage := redis.NewRStorage(cfg, &logger)
	require.NotNil(t, storage)
}

func TestPing_Success(t *testing.T) {
	t.Parallel()

	// Setup miniredis for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	logger := zerolog.New(nil)
	cfg := config.Config{
		Redis: mr.Addr(),
	}

	storage := redis.NewRStorage(cfg, &logger)
	err = storage.Ping(t.Context())
	require.NoError(t, err)
}

func TestPing_Failure(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(nil)
	cfg := config.Config{
		Redis: "invalid-address:6379",
	}

	storage := redis.NewRStorage(cfg, &logger)
	err := storage.Ping(t.Context())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "redis ping")
}

func TestGet_Success(t *testing.T) {
	t.Parallel()

	// Setup miniredis for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Set test data
	testKey := "test-key1"
	testValue := "test-value1"
	err = mr.Set(testKey, testValue)
	require.NoError(t, err)

	logger := zerolog.New(nil)
	cfg := config.Config{
		Redis: mr.Addr(),
	}

	storage := redis.NewRStorage(cfg, &logger)
	value, err := storage.Get(t.Context(), testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)
}

func TestGet_NotFound(t *testing.T) {
	t.Parallel()

	// Setup miniredis for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	logger := zerolog.New(nil)
	cfg := config.Config{
		Redis: mr.Addr(),
	}

	storage := redis.NewRStorage(cfg, &logger)
	_, err = storage.Get(t.Context(), "nonexistent-key")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get value")
}

func TestSet_Success(t *testing.T) {
	t.Parallel()

	// Setup miniredis for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	testKey := "test-key"
	testValue := "test-value"
	expiration := 10 * time.Second

	logger := zerolog.New(nil)
	cfg := config.Config{
		Redis: mr.Addr(),
	}

	storage := redis.NewRStorage(cfg, &logger)
	err = storage.Set(t.Context(), testKey, testValue, expiration)
	require.NoError(t, err)

	// Verify the value was set
	value, err := mr.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)

	// Verify TTL is approximately correct
	ttl := mr.TTL(testKey)
	assert.True(t, ttl <= expiration && ttl > expiration-time.Second)
}

func TestSet_Failure(t *testing.T) {
	t.Parallel()
	// Using redismock to simulate failure
	db, mock := redismock.NewClientMock()
	logger := zerolog.New(nil)
	storage := &redis.RStorage{
		Client: db,
		Logger: &logger,
	}

	testKey := "test-key"
	testValue := "test-value"
	expiration := 10 * time.Second

	// Mock the Set command to fail
	mock.ExpectSet(testKey, testValue, expiration).SetErr(errors.New("redis error"))

	err := storage.Set(t.Context(), testKey, testValue, expiration)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set value")
}

func TestSet_ZeroExpiration(t *testing.T) {
	t.Parallel()
	// Setup miniredis for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	testKey := "test-key"
	testValue := "test-value"

	logger := zerolog.New(nil)
	cfg := config.Config{
		Redis: mr.Addr(),
	}

	storage := redis.NewRStorage(cfg, &logger)
	err = storage.Set(t.Context(), testKey, testValue, 0)
	require.NoError(t, err)

	// Verify the value was set with no expiration
	value, err := mr.Get(testKey)
	require.NoError(t, err)
	assert.Equal(t, testValue, value)

	// Verify there's no TTL
	ttl := mr.TTL(testKey)
	assert.Equal(t, time.Duration(0), ttl)
}

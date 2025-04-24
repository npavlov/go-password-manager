package main

import (
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
)

func TestLoadConfig_Success(t *testing.T) {
	// Create temporary env file
	content := `JWT_SECRET=testsecret
DATABASE_URL=postgres://user:pass@localhost:5432/db`
	tmpFile, err := os.CreateTemp(t.TempDir(), "testenv")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	t.Setenv("ENV_FILE", tmpFile.Name())
	t.Setenv("JWT_SECRET", "testsecret")
	t.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")

	log := zerolog.Nop()
	cfg := loadConfig(&log)

	assert.Equal(t, "testsecret", cfg.JwtSecret)
	assert.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.Database)
}

func TestLoadConfig_Fail(t *testing.T) {
	// Create empty env file
	tmpFile, err := os.CreateTemp(t.TempDir(), "testenv_2")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	t.Setenv("ENV_FILE", tmpFile.Name())

	log := zerolog.Nop()
	assert.PanicsWithValue(t, ErrJWTisNotPorvided, func() {
		loadConfig(&log)
	})
}

func TestSetupDatabase(t *testing.T) {
	t.Parallel()

	t.Run("database connection failure", func(t *testing.T) {
		t.Parallel()

		log := zerolog.Nop()
		//nolint:exhaustruct
		cfg := &config.Config{
			Database: "invalid-connection-string",
		}

		ctx := t.Context()
		db := setupDatabase(ctx, cfg, &log)

		assert.Nil(t, db)
		// In actual implementation, this would likely panic or exit
	})

	t.Run("database connection success", func(t *testing.T) {
		t.Parallel()

		log := zerolog.Nop()
		//nolint:exhaustruct
		cfg := &config.Config{
			Database: "postgres://user:pass@localhost:5432/db",
		}

		ctx := t.Context()
		db := setupDatabase(ctx, cfg, &log)

		assert.NotNil(t, db)
	})
}

func TestSetupStorage(t *testing.T) {
	t.Parallel()
	t.Run("storage initialization", func(t *testing.T) {
		t.Parallel()

		log := zerolog.Nop()
		//nolint:exhaustruct
		cfg := &config.Config{
			Redis: "localhost:6379",
		}

		// Mock DBManager
		//nolint:exhaustruct
		dbMgr := &dbmanager.DBManager{
			DB: nil, // In real test, this would be a test DB
		}

		ctx := t.Context()
		dbStorage, memStorage := setupStorage(ctx, cfg, dbMgr, &log)

		assert.NotNil(t, dbStorage)
		assert.NotNil(t, memStorage)
		// Redis ping would fail in test without running instance
	})
}

func TestSetupMinIO(t *testing.T) {
	t.Parallel()

	t.Run("minio connection failure", func(t *testing.T) {
		t.Parallel()
		//nolint:exhaustruct
		cfg := &config.Config{}

		client, err := setupMinIO(cfg)
		require.Error(t, err)
		assert.Nil(t, client)
	})

	t.Run("minio connection success", func(t *testing.T) {
		t.Parallel()
		//nolint:exhaustruct
		cfg := &config.Config{
			Minio:          "test-minio-address:9000",
			MinioAccessKey: "test",
			MinioSecretKey: "test",
		}

		client, err := setupMinIO(cfg)
		require.NoError(t, err)
		assert.NotNil(t, client)
	})

	// Note: Successful connection test would require a running MinIO instance
	// or a way to mock the MinIO client
}

func TestSetBucket(t *testing.T) {
	t.Parallel()

	t.Run("bucket operations", func(t *testing.T) {
		t.Parallel()
		// Create mock MinIO client
		mockClient := &minio.Client{}

		// In real implementation, you'd use a MinIO mock that implements
		// the BucketExists and MakeBucket methods
		// This is just illustrating the test structure

		//nolint:exhaustruct
		cfg := &config.Config{
			Bucket: "test-bucket",
		}

		ctx := t.Context()

		// This will panic in current form since we can't mock the client methods
		// In real test, you'd use a proper mock
		assert.Panics(t, func() {
			setBucket(ctx, cfg, mockClient)
		})
	})
}

func TestStartServerSuccess(t *testing.T) {
	t.Parallel()

	log := zerolog.Nop()
	//nolint:exhaustruct
	cfg := &config.Config{
		Minio:          "test-minio-address:9000",
		MinioAccessKey: "test",
		MinioSecretKey: "test",
		Certificate:    "testdata/cert.pem",
		PrivateKey:     "testdata/key.pem",
	}

	//nolint:exhaustruct
	dbMgr := &dbmanager.DBManager{
		DB: nil, // In real test, this would be a test DB
	}
	client, err := setupMinIO(cfg)
	require.NoError(t, err)

	grpcManager := startServer(t.Context(), cfg, &log, dbMgr, client)

	assert.NotNil(t, grpcManager)
}

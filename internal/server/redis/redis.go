package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/npavlov/go-password-manager/internal/server/config"
)

type RStorage struct {
	client *redis.Client
	logger *zerolog.Logger
}

func NewRStorage(cfg config.Config, logger *zerolog.Logger) *RStorage {
	// Initialize a Redis client
	//nolint:exhaustruct
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis, // use default Addr
		Password: "",        // no password set
		DB:       0,         // use default DB
	})

	return &RStorage{
		client: redisClient,
		logger: logger,
	}
}

func (rst *RStorage) Ping(ctx context.Context) error {
	// Log the start of the operation
	rst.logger.Info().Msg("Pinging Redis server")

	err := rst.client.Ping(ctx).Err()
	// Record the result in the span
	if err != nil {
		rst.logger.Error().Err(err).Msg("Redis ping failed")

		return errors.Wrap(err, "redis ping")
	}

	rst.logger.Info().Msg("Redis ping successful")

	return nil
}

func (rst *RStorage) Get(ctx context.Context, key string) (string, error) {
	// Log the operation
	rst.logger.Info().Str("key", key).Msg("Getting value from Redis")

	result, err := rst.client.Get(ctx, key).Result()
	if err != nil {
		rst.logger.Error().Err(err).Str("key", key).Msg("Failed to get value from Redis")

		return "", errors.Wrap(err, "failed to get value")
	}

	rst.logger.Info().Str("key", key).Msg("Successfully retrieved value from Redis")

	return result, nil
}

func (rst *RStorage) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	// Log the operation
	rst.logger.Info().
		Str("key", key).
		Dur("expiration", expiration).
		Msg("Setting value in Redis")

	err := rst.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		rst.logger.Error().Err(err).Str("key", key).Msg("Failed to set value in Redis")

		return errors.Wrap(err, "failed to set value")
	}

	rst.logger.Info().Str("key", key).Msg("Successfully set value in Redis")

	return nil
}

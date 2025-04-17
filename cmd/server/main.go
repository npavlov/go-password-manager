package main

import (
	"context"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/npavlov/go-password-manager/internal/server/adapter"
	"github.com/npavlov/go-password-manager/internal/server/buildinfo"
	"github.com/npavlov/go-password-manager/internal/server/config"
	"github.com/npavlov/go-password-manager/internal/server/dbmanager"
	"github.com/npavlov/go-password-manager/internal/server/redis"
	"github.com/npavlov/go-password-manager/internal/server/service"
	"github.com/npavlov/go-password-manager/internal/server/service/auth"
	"github.com/npavlov/go-password-manager/internal/server/service/card"
	"github.com/npavlov/go-password-manager/internal/server/service/file"
	"github.com/npavlov/go-password-manager/internal/server/service/item"
	"github.com/npavlov/go-password-manager/internal/server/service/meta"
	"github.com/npavlov/go-password-manager/internal/server/service/note"
	"github.com/npavlov/go-password-manager/internal/server/service/password"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/npavlov/go-password-manager/internal/utils"
)

var (
	ErrDatabaseNotConnected = errors.New("database is not connected")
	ErrJWTisNotPorvided     = errors.New("JWT token is not provided")
	ErrMinioNotConnected    = errors.New("MinIO is not connected")
)

func main() {
	log := logger.NewLogger(zerolog.DebugLevel).Get()

	log.Info().Str("buildVersion", buildinfo.Version).
		Str("buildCommit", buildinfo.Commit).
		Str("buildDate", buildinfo.Date).Msg("Starting server")

	cfg := loadConfig(&log)
	var wg sync.WaitGroup

	ctx, cancel := utils.WithSignalCancel(context.Background(), &log)
	dbManager := setupDatabase(ctx, cfg, &log)

	if dbManager == nil {
		cancel()
		dbManager.Close()
		log.Fatal().Err(ErrDatabaseNotConnected).Msg("Database is not connected")
	}

	defer cancel()
	defer dbManager.Close()

	wg.Add(1)
	starServer(ctx, cfg, &log, &wg, dbManager)

	utils.WaitForShutdown(&wg)
}

func loadConfig(log *zerolog.Logger) *config.Config {
	if err := godotenv.Load("server.env"); err != nil {
		log.Error().Err(err).Msg("Error loading server.env file")
	}

	cfg := config.NewConfigBuilder(log).FromEnv().FromFlags().Build()
	log.Info().Interface("config", cfg).Msg("Configuration loaded")

	if cfg.JwtSecret == "" {
		panic(ErrJWTisNotPorvided)
	}

	return cfg
}

func starServer(
	ctx context.Context,
	cfg *config.Config,
	log *zerolog.Logger,
	wg *sync.WaitGroup,
	dbM *dbmanager.DBManager,
) {
	dbStorage, memStorage := setupStorage(ctx, cfg, dbM, log)

	//nolint:contextcheck
	grpcManager := service.NewGRPCManager(cfg, log, memStorage)
	grpcServer := grpcManager.GetServer()

	authService := auth.NewAuthService(log, dbStorage, cfg, memStorage)
	authService.RegisterService(grpcServer)

	passwordService := password.NewPasswordService(log, dbStorage, cfg)
	passwordService.RegisterService(grpcServer)

	noteService := note.NewNoteService(log, dbStorage, cfg)
	noteService.RegisterService(grpcServer)

	cardService := card.NewCardService(log, dbStorage, cfg)
	cardService.RegisterService(grpcServer)

	minioClient, err := setupMinIO(cfg)
	if err != nil {
		log.Fatal().Err(ErrMinioNotConnected).Msg("Failed to connect to MinIO")
	}

	setBucket(ctx, cfg, minioClient)

	fileService := file.NewFileService(log, dbStorage, cfg, adapter.NewMinioAdapter(minioClient))
	fileService.RegisterService(grpcServer)

	itemService := item.NewItemService(log, dbStorage, cfg)
	itemService.RegisterService(grpcServer)

	metaService := meta.NewMetadataService(log, dbStorage, cfg)
	metaService.RegisterService(grpcServer)

	grpcManager.Start(ctx, wg)
}

func setupDatabase(ctx context.Context, cfg *config.Config, log *zerolog.Logger) *dbmanager.DBManager {
	dbManager := dbmanager.NewDBManager(cfg.Database, log).Connect(ctx)
	if dbManager.DB == nil {
		return nil
	}
	dbManager.VerifyConnection(ctx).ApplyMigrations()

	return dbManager
}

//nolint:ireturn
func setupStorage(
	ctx context.Context,
	cfg *config.Config,
	dbManager *dbmanager.DBManager,
	log *zerolog.Logger,
) (*storage.DBStorage, redis.MemStorage) {
	st := storage.NewDBStorage(dbManager.DB, log)
	memStorage := redis.NewRStorage(*cfg, log)

	if err := memStorage.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Error connecting to redis")
	}

	return st, memStorage
}

func setupMinIO(cfg *config.Config) (*minio.Client, error) {
	//nolint:exhaustruct
	client, err := minio.New(cfg.Minio, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: false, // Set to `true` if using HTTPS
	})

	return client, errors.Wrap(err, "error creating minio client")
}

func setBucket(ctx context.Context, cfg *config.Config, client *minio.Client) {
	// Check if bucket exists, create if not
	bucketName := cfg.Bucket
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to check MinIO bucket")
	}
	if !exists {
		//nolint:exhaustruct
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatal().Err(err).Msg("Failed to create MinIO bucket")
		}
		log.Info().Str("bucket", bucketName).Msg("MinIO bucket created")
	} else {
		log.Info().Str("bucket", bucketName).Msg("MinIO bucket already exists")
	}

	log.Info().Msg("Connected to MinIO successfully")
}

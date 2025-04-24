package config

import (
	"flag"
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"

	"github.com/npavlov/go-password-manager/internal/utils"
)

type Config struct {
	Address          string `env:"ADDRESS"          envDefault:":9090"`
	Database         string `env:"DATABASE_DSN"     envDefault:""`
	JwtSecret        string `env:"JWT_SECRET"       envDefault:""`
	Certificate      string `env:"CERTIFICATE"      envDefault:""`
	PrivateKey       string `env:"PRIVATE_KEY"      envDefault:""`
	MasterKey        string `env:"MASTER_KEY"       envDefault:""`
	Redis            string `env:"REDIS"            envDefault:"localhost:6379"`
	Minio            string `env:"MINIO"            envDefault:""`
	MinioAccessKey   string `env:"MINIO_ACCESS_KEY" envDefault:""`
	MinioSecretKey   string `env:"MINIO_SECRET_KEY" envDefault:""`
	Bucket           string `env:"BUCKET"           envDefault:"encrypted-bucket"`
	SecuredMasterKey utils.ISecureString
}

// Builder defines the builder for the Config struct.
type Builder struct {
	cfg    *Config
	logger *zerolog.Logger
	mu     sync.RWMutex
}

// NewConfigBuilder initializes the ConfigBuilder with default values.
func NewConfigBuilder(log *zerolog.Logger) *Builder {
	return &Builder{
		cfg: &Config{
			Address:          "",
			Database:         "",
			JwtSecret:        "",
			Certificate:      "",
			PrivateKey:       "",
			Redis:            "",
			MasterKey:        "",
			Bucket:           "",
			Minio:            "",
			MinioAccessKey:   "",
			MinioSecretKey:   "",
			SecuredMasterKey: nil,
		},
		logger: log,
		mu:     sync.RWMutex{},
	}
}

// FromEnv parses environment variables into the ConfigBuilder.
func (b *Builder) FromEnv() *Builder {
	if err := env.Parse(b.cfg); err != nil {
		b.logger.Error().Err(err).Msg("failed to parse environment variables")
	}

	return b
}

// FromFlags parses command line flags into the ConfigBuilder.
func (b *Builder) FromFlags() *Builder {
	b.mu.Lock()
	defer b.mu.Unlock()
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.StringVar(&b.cfg.Address, "a", b.cfg.Address, "address and port to run server")
	fs.StringVar(&b.cfg.Database, "d", b.cfg.Database, "database DSN")
	fs.StringVar(&b.cfg.JwtSecret, "jwt", b.cfg.JwtSecret, "JWT Secret")
	fs.StringVar(&b.cfg.Certificate, "cert", b.cfg.Certificate, "Certificate")
	fs.StringVar(&b.cfg.PrivateKey, "privatekey", b.cfg.PrivateKey, "Private Key for http connection")
	fs.StringVar(&b.cfg.Redis, "redis", b.cfg.Redis, "Redis connection string")
	fs.StringVar(&b.cfg.MasterKey, "masterkey", b.cfg.MasterKey, "Master Key for encrypting data")
	fs.StringVar(&b.cfg.Minio, "minio", b.cfg.Minio, "Minio address")
	fs.StringVar(&b.cfg.MinioAccessKey, "minio_access_key", b.cfg.MinioAccessKey, "Minio access key")
	fs.StringVar(&b.cfg.MinioSecretKey, "minio_secret_key", b.cfg.MinioSecretKey, "Minio secret key")
	fs.StringVar(&b.cfg.Bucket, "bucket", b.cfg.Bucket, "Bucket name for Minio")
	_ = fs.Parse(os.Args[1:])

	return b
}

// FromObj sets cfg from object.
func (b *Builder) FromObj(cfg *Config) *Builder {
	b.cfg = cfg

	return b
}

// Build returns the final configuration.
func (b *Builder) Build() *Config {
	b.cfg.SecuredMasterKey = utils.NewString(b.cfg.MasterKey)
	b.cfg.MasterKey = ""

	return b.cfg
}

//nolint:tagliatelle
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
)

type Config struct {
	Address     string `env:"ADDRESS"        envDefault:":9090"`
	Database    string `env:"DATABASE_DSN"          envDefault:""`
	JwtSecret   string `env:"JWT_SECRET"             envDefault:""`
	Certificate string `env:"CERTIFICATE"          envDefault:""`
	PrivateKey  string `env:"PRIVATE_KEY"          envDefault:""`
	Redis       string `env:"REDIS"                  envDefault:"localhost:6379"`
}

// Builder defines the builder for the Config struct.
type Builder struct {
	cfg    *Config
	logger *zerolog.Logger
}

// NewConfigBuilder initializes the ConfigBuilder with default values.
func NewConfigBuilder(log *zerolog.Logger) *Builder {
	return &Builder{
		cfg: &Config{
			Address:     "",
			Database:    "",
			JwtSecret:   "",
			Certificate: "",
			PrivateKey:  "",
			Redis:       "",
		},
		logger: log,
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
	flag.StringVar(&b.cfg.Address, "a", b.cfg.Address, "address and port to run server")
	flag.StringVar(&b.cfg.Database, "d", b.cfg.Database, "database DSN")
	flag.StringVar(&b.cfg.JwtSecret, "jwt", b.cfg.JwtSecret, "JWT Secret")
	flag.StringVar(&b.cfg.Certificate, "cert", b.cfg.Certificate, "Certificate")
	flag.StringVar(&b.cfg.PrivateKey, "privatekey", b.cfg.PrivateKey, "Private Key")
	flag.StringVar(&b.cfg.Redis, "redis", b.cfg.Redis, "Redis connection string")
	flag.Parse()

	return b
}

// FromObj sets cfg from object.
func (b *Builder) FromObj(cfg *Config) *Builder {
	b.cfg = cfg

	return b
}

// Build returns the final configuration.
func (b *Builder) Build() *Config {
	return b.cfg
}

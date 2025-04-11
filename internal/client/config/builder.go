package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"

	"github.com/npavlov/go-password-manager/internal/utils"
)

type Config struct {
	Address          string `env:"ADDRESS"     envDefault:":9090"`
	MasterKey        string `env:"MASTER_KEY"  envDefault:""`
	Certificate      string `env:"CERTIFICATE" envDefault:""`
	TokenFile        string `env:"TOKEN_FILE"  envDefault:""`
	SecuredMasterKey utils.ISecureString
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
			Address:          "",
			MasterKey:        "",
			Certificate:      "",
			TokenFile:        "",
			SecuredMasterKey: nil,
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
	flag.StringVar(&b.cfg.MasterKey, "masterkey", b.cfg.MasterKey, "Master Key for encrypting data")
	flag.StringVar(&b.cfg.Certificate, "cert", b.cfg.Certificate, "Certificate")
	flag.StringVar(&b.cfg.TokenFile, "token_file", b.cfg.TokenFile, "File where do we store tokens")
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
	b.cfg.SecuredMasterKey = utils.NewString(b.cfg.MasterKey)
	b.cfg.MasterKey = ""

	return b.cfg
}

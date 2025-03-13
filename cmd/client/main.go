package client

import (
	"fmt"

	"github.com/npavlov/go-password-manager/internal/client/buildinfo"
	"github.com/npavlov/go-password-manager/internal/pkg/logger"
	"github.com/rs/zerolog"
)

func main() {
	log := logger.NewLogger(zerolog.DebugLevel).Get()

	log.Info().Str("buildVersion", buildinfo.Version).
		Str("buildCommit", buildinfo.Commit).
		Str("buildDate", buildinfo.Date).Msg("Starting agent")

	fmt.Println("Hello World")
}

package migrate

import (
	"time"
	"turbo-mailer-server/internal/query"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("migrate database")
		start := time.Now()
		if err := query.Migrate(); err != nil {
			log.Fatal().Err(err).Msg("migrate failed")
		}
		log.Info().Dur("duration", time.Since(start)).Msg("migrate database success")
	},
}

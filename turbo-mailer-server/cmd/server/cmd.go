package server

import (
	_ "turbo-mailer-server/docs"
	"turbo-mailer-server/internal/dispatch"
	"turbo-mailer-server/internal/worker"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	mode string
)

var Cmd = &cobra.Command{
	Use:   "serve",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		// initialize.Do(cmd.Context())

		// TODO 还是需要放到一起，因为其他的模式也需要 web 服务，debug，监控等

		switch mode {
		case "api":
			api()
		case "worker":
			worker.Boot(cmd.Context())
		case "dispatcher":
			dispatch.Boot(cmd.Context())
		default:
			log.Fatal().Str("mode", mode).Msg("Invalid mode")
		}

	},
}

func init() {
	Cmd.Flags().StringVarP(&mode, "mode", "m", "api", "Run as server mode (api, worker, dispatcher)")
}

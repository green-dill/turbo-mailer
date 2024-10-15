package cmd

import (
	"fmt"
	"os"
	"turbo-mailer-server/cmd/migrate"
	"turbo-mailer-server/cmd/server"
	"turbo-mailer-server/internal/config"
	"turbo-mailer-server/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "turbo-mailer-server",
	SilenceUsage: true,
	Version:      fmt.Sprintf("%s-%s", version.Version, version.CommitID),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.Initialize()
	},
}

// add sub command
func init() {
	rootCmd.AddCommand(server.Cmd)
	rootCmd.AddCommand(migrate.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

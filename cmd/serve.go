package cmd

import (
	"github.com/IceSandwich/IceVoice/application"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func serveRun(cmd *cobra.Command, args []string) {
	if err := application.Init(8888); err != nil {
		log.WithError(err).Fatalf("cannot init application")
	}

	application.Run()
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run as a server",
	Run:   serveRun,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

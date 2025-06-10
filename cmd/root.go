package cmd

import (
	"fmt"
	"os"

	"github.com/inconshreveable/mousetrap"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "IceVoice",
	Short: "IceVoice is an ollama-like voice application",
	Long:  `IceVoice is a voice application served as backend app, includes features like text-to-speech, speech-to-text, speech-recognition and etc.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if mousetrap.StartedByExplorer() && len(os.Args) == 1 {
		fmt.Printf("Starting GUI...\n")
		os.Args = append(os.Args, "gui")
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Disable cobra mousetrap because we use custom mousetrap to show gui when args is not provided.
	cobra.MousetrapHelpText = ""

	// Hide completion options from help text
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

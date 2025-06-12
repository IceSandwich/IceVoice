package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/IceSandwich/IceVoice/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func rmRun(cmd *cobra.Command, args []string) {
	freq := make(map[string]int)

	var targetPath string
	var targetConf config.ModelConfig

	config.WalkManifest(func(path string, model config.ModelConfig) error {
		for _, f := range model.Rootfs {
			if _, ok := freq[f]; !ok {
				freq[f] = 0
			}
			freq[f]++
		}

		modelID := fmt.Sprintf("%s:%s", model.Name, model.ModelType)
		if args[0] == modelID {
			targetPath = path
			targetConf = model
		}

		return nil
	})

	if targetPath == "" {
		log.Errorf("No found model %s", args[0])
		return
	}

	if err := os.Remove(targetPath); err != nil {
		log.WithError(err).Fatalf("Failed to remove config file %s", targetPath)
	}

	if config.IsManifestEmpty(targetConf.Name) {
		baseDir := filepath.Dir(targetPath)
		if err := os.Remove(baseDir); err != nil {
			log.WithError(err).Errorf("Failed to remove base dir %s", baseDir)
		}
	}

	for _, hash := range targetConf.Rootfs {
		filename := config.HashToBlobFile(hash)
		if _, err := os.Stat(filename); err != nil {
			// file already removed
			freq[hash] = 0
		} else {
			freq[hash]--
		}
	}

	for hash, count := range freq {
		if count == 0 {
			filename := config.HashToBlobFile(hash)
			if err := os.Remove(filename); err != nil {
				log.Errorf("Failed to remove %s (%v)", filename, err)
			}
		}
	}

	log.Infof("Successfully removed %s", args[0])
}

var rmCmd = &cobra.Command{
	Use:   "rm modelfile",
	Short: "Remove a model installed",
	Args:  cobra.MinimumNArgs(1),
	Run:   rmRun,
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

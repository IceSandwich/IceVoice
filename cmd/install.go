/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/IceSandwich/IceVoice/config"
	"github.com/IceSandwich/IceVoice/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func installRun(cmd *cobra.Command, args []string) {
	filename := args[0]

	zfile, err := OpenZip(filename)
	if err != nil {
		log.WithError(err).Fatalf("failed to open %s", filename)
	}
	defer zfile.Close()

	if !zfile.Exists("config.json") {
		log.Fatal("failed to find config.json file in zip archive")
	}

	conf := config.ModelConfig{}
	if err := zfile.ReadJson("config.json", &conf); err != nil {
		log.WithError(err).Fatal("failed to read config file in zip archive")
	}

	model := models.Create(conf.Architecture)
	if model == nil {
		log.Fatalf("failed to create model for %s", conf.Architecture)
	}

	requires := model.GetRequireFiles()
	for _, filename := range requires {
		actualFilename, ok := conf.Rootfs[filename]
		if !ok {
			log.Fatalf("%s is required by %s but not found in zip archive", filename, conf.Name)
		}

		if !zfile.Exists(actualFilename) {
			log.Fatalf("file %s not found in zip archive which is the actual file of %s", actualFilename, filename)
		}
	}

	resdir := filepath.Join(config.GetManifestsDir(), conf.Name)
	if err := os.MkdirAll(resdir, os.ModePerm); err != nil {
		log.WithError(err).Fatalf("failed to create model directory %s", resdir)
	}

	newConf := config.ModelConfig{}
	newConf.Name = conf.Name
	newConf.Task = conf.Task
	newConf.Architecture = conf.Architecture
	newConf.ModelType = conf.ModelType
	newConf.Quantization = conf.Quantization
	newConf.Rootfs = make(map[string]string)

	for _, filename := range requires {
		actualFilename := conf.Rootfs[filename]

		hash256 := zfile.CalcHash256(actualFilename)
		if hash256 == "" {
			log.Fatalf("failed to calculate hash of file %s in zip archive", filename)
		}

		targetFilename := filepath.Join(config.GetBlobsDir(), hash256)
		if err := zfile.Extract(actualFilename, targetFilename); err != nil {
			log.WithError(err).Fatalf("failed to extract %s from zip archive", filename)
		}

		newConf.Rootfs[filename] = fmt.Sprintf("sha256:%s", hash256)
	}

	configFilename := filepath.Join(resdir, conf.ModelType+".json")
	data, err := json.MarshalIndent(&newConf, "", "\t")
	if err != nil {
		log.WithError(err).Fatalf("failed to marshal model config")
	}
	if err := os.WriteFile(configFilename, data, os.ModePerm); err != nil {
		log.WithError(err).Fatalf("failed to write model config %s", configFilename)
	}

	log.Infof("Successfully install model %s", newConf.Name)
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install modelfile",
	Args:  cobra.MinimumNArgs(1),
	Short: "Install a model from zip file",
	Run:   installRun,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

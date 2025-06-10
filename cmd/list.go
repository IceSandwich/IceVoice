package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/IceSandwich/IceVoice/config"
	"github.com/mattn/go-runewidth"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	TaskTTS = "tts"
)

type modelBag struct {
	TTS []config.ModelConfig
}

func (m *modelBag) ShowTTS() {
	data := [][]string{
		{"NAME", "ARCHITECTURE", "SIZE"},
	}

	for _, m := range m.TTS {
		size := uint64(0)
		skip := false
		for _, hash := range m.Rootfs {
			filename := config.HashToBlobFile(hash)
			finfo, err := os.Stat(filename)
			if err != nil {
				log.WithError(err).Errorf("broken file %s for model %s", filename, m.Name)
				skip = true
				break
			}

			size += uint64(finfo.Size())
		}
		if skip {
			continue
		}

		data = append(data, []string{
			fmt.Sprintf("%s:%s", m.Name, m.ModelType),
			m.Architecture,
			HumanReadableBytes(size),
		})
	}

	fmt.Println("Available Text-to-speech models:")
	table := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
		Symbols: nil,
	})))
	table.Header(data[0])
	table.Bulk(data[1:])
	table.Render()
}

func listRun(cmd *cobra.Command, args []string) {
	bag := modelBag{}

	filepath.Walk(config.GetManifestsDir(), func(path string, info os.FileInfo, ex error) error {
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			log.WithError(err).Fatalf("failed to open manifest file %s, error: %v+", path, err)
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			log.WithError(err).Fatalf("failed to read manifest file %s, error: %v+", path, err)
		}

		var model config.ModelConfig
		if err := json.Unmarshal(data, &model); err != nil {
			log.WithError(err).Fatalf("failed to parse manifest file %s, error: %v+", path, err)
		}

		switch model.Task {
		case TaskTTS:
			bag.TTS = append(bag.TTS, model)
		default:
			log.Fatalf("unknown model task %s in manifest file %s", model.Task, path)
		}

		return nil
	})

	bag.ShowTTS()
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all models installed",
	Run:   listRun,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// v1.0.7 is the latest version in gopkg so we add these lines to fix the issue
	// https://github.com/olekukonko/tablewriter/issues/275
	// TODO: remove this when tablewriter >= v1.0.8
	runewidth.EastAsianWidth = false
	runewidth.DefaultCondition.EastAsianWidth = false
}

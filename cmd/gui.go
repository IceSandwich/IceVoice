package cmd

import (
	"strings"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/spf13/cobra"
)

func guiRun(cmd *cobra.Command, args []string) {
	var inTE, outTE *walk.TextEdit

	declarative.MainWindow{
		Title: "SCREAMO",
		MinSize: declarative.Size{
			Width:  600,
			Height: 400,
		},
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.HSplitter{
				Children: []declarative.Widget{
					declarative.TextEdit{AssignTo: &inTE},
					declarative.TextEdit{AssignTo: &outTE, ReadOnly: true},
				},
			},
			declarative.PushButton{
				Text: "SCREAM",
				OnClicked: func() {
					outTE.SetText(strings.ToUpper(inTE.Text()))
				},
			},
		},
	}.Run()
}

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Open a simple GUI to control the IceVoice",
	Run:   guiRun,
}

func init() {
	rootCmd.AddCommand(guiCmd)
}

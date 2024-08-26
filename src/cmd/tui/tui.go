package tui

import (
	"os"

	"github.com/defenseunicorns/lula/src/internal/tui"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

type flags struct {
	InputFile string // -f --input-file
}

var opts = &flags{}

var tuiHelp = `
To view an OSCAL model in the TUI:
	lula tui -f /path/to/oscal-component.yaml
`

var tuiCmd = &cobra.Command{
	Use:     "tui",
	Short:   "TUI viewer for OSCAL models",
	Long:    "TUI viewer for OSCAL models",
	Example: tuiHelp,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the OSCAL model from the file
		data, err := os.ReadFile(opts.InputFile)
		if err != nil {
			message.Fatalf(err, "error reading file: %v", err)
		}
		oscalModel, err := oscal.NewOscalModel(data)
		if err != nil {
			message.Fatalf(err, "error creating oscal model from file: %v", err)
		}

		// Add debugging
		// TODO: need to integrate with the log file handled by messages
		if message.GetLogLevel() == message.DebugLevel {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				message.Fatalf(err, err.Error())
			}
			defer f.Close()
		}

		p := tea.NewProgram(tui.NewOSCALModel(*oscalModel), tea.WithAltScreen(), tea.WithMouseCellMotion())

		if _, err := p.Run(); err != nil {
			message.Fatalf(err, err.Error())
		}
	},
}

func TuiCommand() *cobra.Command {
	tuiCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to the target OSCAL model")
	tuiCmd.MarkFlagRequired("input-file")
	return tuiCmd
}

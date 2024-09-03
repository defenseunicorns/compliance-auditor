package console

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

var consoleHelp = `
To view an OSCAL model in the Console:
	lula console -f /path/to/oscal-component.yaml
`

var consoleLong = `
The Lula Console is a text-based terminal user interface that allows users to 
interact with the OSCAL documents in a more intuitive and visual way.
`

var consoleCmd = &cobra.Command{
	Use:     "console",
	Aliases: []string{"ui"},
	Short:   "Console terminal user interface for OSCAL models",
	Long:    consoleLong,
	Example: consoleHelp,
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

func ConsoleCommand() *cobra.Command {
	consoleCmd.Flags().StringVarP(&opts.InputFile, "input-file", "f", "", "the path to the target OSCAL model")
	consoleCmd.MarkFlagRequired("input-file")
	return consoleCmd
}

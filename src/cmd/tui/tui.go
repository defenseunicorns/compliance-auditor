package tui

import (
	"fmt"
	"os"

	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/internal/tui"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

var tuiHelp = `
To view an OSCAL model in the TUI:
	lula tui /path/to/oscal-component.yaml
`

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "TUI viewer for OSCAL models",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SkipLogFile = true
	},
	Long:    "TUI viewer for OSCAL models",
	Example: tuiHelp,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			message.Fatal(fmt.Errorf("no input file provided"), "no input file provided")
		}

		// Get the OSCAL model from the file
		data, err := os.ReadFile(args[0])
		if err != nil {
			message.Fatalf(err, "error reading file: %v", err)
		}
		oscalModel, err := oscal.NewOscalModel(data)
		if err != nil {
			message.Fatalf(err, "error creating oscal model from file: %v", err)
		}

		p := tea.NewProgram(tui.NewOSCALModel(*oscalModel), tea.WithAltScreen(), tea.WithMouseCellMotion())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Include(rootCmd *cobra.Command) {
	rootCmd.AddCommand(tuiCmd)
}

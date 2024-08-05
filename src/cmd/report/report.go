package report

import (
	"github.com/spf13/cobra"
)

type flags struct {
	InputFile  string // -f --input-file
	OutputFile string // -o --output-file
}

var opts = &flags{}

var reportHelp = `

`

var reportCmd = &cobra.Command{
	Use:     "report",
	Hidden:  false,
	Aliases: []string{"r"},
	Short:   "Build a compliance report",
	Example: reportHelp,
	Run: func(_ *cobra.Command, args []string) {

	},
}

func ReportCommand() *cobra.Command {

	reportFlags()

	return reportCmd
}

func reportFlags() {
	reportFlags := reportCmd.PersistentFlags()

	reportFlags.StringVarP(&opts.InputFile, "input-file", "f", "", "Path to a manifest file")
	reportFlags.StringVarP(&opts.OutputFile, "output-file", "o", "", "Path and Name to an output file")

}

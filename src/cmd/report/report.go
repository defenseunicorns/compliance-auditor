package report

import (
	"fmt"

	"github.com/defenseunicorns/lula/src/internal/reporting"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
)

var reportHelp = `
To create a new report:
lula report -f oscal-component-definition.yaml

To create a new report in json format:
lula report -f oscal-component-definition.yaml --file-format json

To create a new report in yaml format:
lula report -f oscal-component-definition.yaml --file-format yaml
`

func ReportCommand() *cobra.Command {
	var (
		inputFile  string
		fileFormat string
	)

	cmd := &cobra.Command{
		Use:     "report",
		Short:   "Build a compliance report",
		Example: reportHelp, // reuse your existing help text
		RunE: func(cmd *cobra.Command, args []string) error {
			err := reporting.GenerateReport(inputFile, fileFormat)
			if err != nil {
				return fmt.Errorf("error generating report: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "Path to an OSCAL file")
	cmd.Flags().StringVar(&fileFormat, "file-format", "table", "File format of the report")
	err := cmd.MarkFlagRequired("input-file")
	if err != nil {
		message.Fatal(err, "error initializing report command flags")
	}

	return cmd
}

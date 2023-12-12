package cmd

import (
	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/cmd/evaluate"
	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/cmd/version"
)

var rootCmd = &cobra.Command{
	Use:   "lula",
	Short: "Risk Management as Code",
	Long:  `Real Time Risk Transparency through automated validation`,
}

func Execute() {

	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	commands := []*cobra.Command{
		validate.ValidateCommand(),
		evaluate.EvaluateCommand(),
	}

	rootCmd.AddCommand(commands...)
	tools.Include(rootCmd)
	version.Include(rootCmd)
}

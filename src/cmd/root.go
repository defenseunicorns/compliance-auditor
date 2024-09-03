package cmd

import (
	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/cmd/console"
	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/cmd/evaluate"
	"github.com/defenseunicorns/lula/src/cmd/generate"
	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/cmd/version"
)

var LogLevelCLI string

var rootCmd = &cobra.Command{
	Use: "lula",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		common.SetupClI(LogLevelCLI)
	},
	Short: "Risk Management as Code",
	Long:  `Real Time Risk Transparency through automated validation`,
}

func Execute() {

	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	v := common.InitViper()

	commands := []*cobra.Command{
		validate.ValidateCommand(),
		evaluate.EvaluateCommand(),
		generate.GenerateCommand(),
		console.ConsoleCommand(),
	}

	rootCmd.AddCommand(commands...)
	tools.Include(rootCmd)
	version.Include(rootCmd)
	dev.Include(rootCmd)

	rootCmd.PersistentFlags().StringVarP(&LogLevelCLI, "log-level", "l", v.GetString(common.VLogLevel), "Log level when running Lula. Valid options are: warn, info, debug, trace")
}

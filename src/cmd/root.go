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

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "lula",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			common.SetupClI(LogLevelCLI)
		},
		Short: "Risk Management as Code",
		Long:  `Real Time Risk Transparency through automated validation`,
	}

	v := common.InitViper()

	commands := []*cobra.Command{
		validate.ValidateCommand(),
		evaluate.EvaluateCommand(),
		generate.GenerateCommand(),
		console.ConsoleCommand(),
	}

	cmd.AddCommand(commands...)
	tools.Include(cmd)
	version.Include(cmd)
	dev.Include(cmd)
	cmd.AddCommand(internalCmd)

	cmd.PersistentFlags().StringVarP(&LogLevelCLI, "log-level", "l", v.GetString(common.VLogLevel), "Log level when running Lula. Valid options are: warn, info, debug, trace")

	return cmd
}

func RootCommand() *cobra.Command {
	return newRootCmd()
}

func Execute() {
	cobra.CheckErr(newRootCmd().Execute())
}

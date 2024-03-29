package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/cmd/evaluate"
	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/cmd/version"
	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

var LogLevelCLI string

var rootCmd = &cobra.Command{
	Use: "lula",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		match := map[string]message.LogLevel{
			"warn":  message.WarnLevel,
			"info":  message.InfoLevel,
			"debug": message.DebugLevel,
			"trace": message.TraceLevel,
		}

		// No log level set, so use the default
		if LogLevelCLI != "" {
			if lvl, ok := match[LogLevelCLI]; ok {
				message.SetLogLevel(lvl)
				message.Debug("Log level set to " + LogLevelCLI)
			} else {
				message.Warn("Invalid log level. Valid options are: warn, info, debug, trace.")
			}
		}

		// Disable progress bars for CI envs
		if os.Getenv("CI") == "true" {
			message.Debug("CI environment detected, disabling progress bars")
			message.NoProgress = true
		}

		if !config.SkipLogFile {
			message.UseLogFile()
		}
	},
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
	dev.Include(rootCmd)

	rootCmd.PersistentFlags().StringVarP(&LogLevelCLI, "log-level", "l", "info", "Log level when running Lula. Valid options are: warn, info, debug, trace")
}

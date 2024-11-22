package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/defenseunicorns/lula/src/cmd/console"
	"github.com/defenseunicorns/lula/src/cmd/dev"
	"github.com/defenseunicorns/lula/src/cmd/evaluate"
	"github.com/defenseunicorns/lula/src/cmd/generate"
	"github.com/defenseunicorns/lula/src/cmd/report"
	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/cmd/version"
	"github.com/spf13/cobra"
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

func RootCommand() *cobra.Command {

	cmd := rootCmd

	return cmd
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		select {
		case <-c:
			fmt.Println("Got signal, shutting down...")
			cancel()
			os.Exit(2)
		case <-ctx.Done():
			return
		}
	}()

	cobra.CheckErr(rootCmd.ExecuteContext(ctx))
}

func init() {

	v := common.InitViper()

	commands := []*cobra.Command{
		validate.ValidateCommand(),
		evaluate.EvaluateCommand(),
		generate.GenerateCommand(),
		report.ReportCommand(),
		console.ConsoleCommand(),
		dev.DevCommand(),
	}

	rootCmd.AddCommand(commands...)
	tools.Include(rootCmd)
	version.Include(rootCmd)

	rootCmd.PersistentFlags().StringVarP(&LogLevelCLI, "log-level", "l", v.GetString(common.VLogLevel), "Log level when running Lula. Valid options are: warn, info, debug, trace")
}

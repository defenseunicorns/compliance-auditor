package tools

import (
	"github.com/defenseunicorns/lula/src/config"
	"github.com/spf13/cobra"
)

var toolsCmd = &cobra.Command{
	Use:     "tools",
	Aliases: []string{"t"},
	Short:   "Collection of additional commands to make OSCAL easier",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.SkipLogFile = true
		// Call the parent's (root) PersistentPreRun
		if parentPreRun := cmd.Parent().PersistentPreRun; parentPreRun != nil {
			parentPreRun(cmd.Parent(), args)
		}
	},
}

// Include adds the tools command to the root command.
func Include(rootCmd *cobra.Command) {
	rootCmd.AddCommand(toolsCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate compliance artifact templates",
	Long:  `Generation of compliance artifact templates`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")
	},
}

func GenerateCommand() *cobra.Command {

	return generateCmd
}

package cmd

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"

	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var internalCmd = &cobra.Command{
	Use:    "internal",
	Hidden: true,
	Short:  "Internal Commands",
}

var genCLIDocs = &cobra.Command{
	Use:   "gen-cli-docs",
	Short: "Generate CLI command documentation",
	RunE: func(_ *cobra.Command, _ []string) error {
		// Don't include the datestamp in the output
		rootCmd.DisableAutoGenTag = true

		// rootCmd.RemoveCommand()

		// remove existing docs
		glob, err := filepath.Glob("./docs/cli-commands/lula*.md")
		if err != nil {
			return err
		}
		for _, f := range glob {
			err := os.Remove(f)
			if err != nil {
				return err
			}
		}

		var prependTitle = func(s string) string {
			fmt.Println(s)
			name := filepath.Base(s)

			// strip .md extension
			name = name[:len(name)-3]

			// replace _ with space
			title := strings.Replace(name, "_", " ", -1)

			return fmt.Sprintf(`---
title: %s
description: Lula CLI command reference for <code>%s</code>.
type: docs
---
`, title, title)
		}

		var linkHandler = func(link string) string {
			return "./" + link
		}

		err = doc.GenMarkdownTreeCustom(rootCmd, "./docs/cli-commands", prependTitle, linkHandler)
		if err != nil {
			return err
		}

		message.Success("Internal documentation generated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(internalCmd)

	internalCmd.AddCommand(genCLIDocs)
}

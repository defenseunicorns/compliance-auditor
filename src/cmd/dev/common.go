package dev

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	cmdCommon "github.com/defenseunicorns/lula/src/cmd/common"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

const STDIN = "0"
const NO_TIMEOUT = -1
const DEFAULT_TIMEOUT = 1

func DevCommand() *cobra.Command {
	var (
		setOpts []string
	)

	cmd := &cobra.Command{
		Use:     "dev",
		Aliases: []string{"d"},
		Short:   "Collection of dev commands to make dev life easier",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.SkipLogFile = true
			// Call the parent's (root) PersistentPreRun
			if parentPreRun := cmd.Parent().PersistentPreRun; parentPreRun != nil {
				parentPreRun(cmd.Parent(), args)
			}
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&setOpts, "set", "s", []string{}, "set a value in the template data")

	cmd.AddCommand(DevLintCommand())
	cmd.AddCommand(DevValidateCommand())
	cmd.AddCommand(DevGetResourcesCommand())

	return cmd
}

var RunInteractively bool = true // default to run dev command interactively

// ReadValidation reads the validation yaml file and returns the validation bytes
func ReadValidation(cmd *cobra.Command, spinner *message.Spinner, path string, timeout int) ([]byte, error) {
	var validationBytes []byte
	var err error

	if path == STDIN {
		var inputReader io.Reader = cmd.InOrStdin()

		// If the timeout is not -1, wait for the timeout then close and return an error
		go func() {
			if timeout != NO_TIMEOUT {
				time.Sleep(time.Duration(timeout) * time.Second)
				//nolint:errcheck
				cmd.Help() // #nosec G104
				message.Fatalf(fmt.Errorf("timed out waiting for stdin"), "timed out waiting for stdin")
			}
		}()

		// Update the spinner message
		spinner.Updatef("reading from stdin...")
		// Read from stdin
		validationBytes, err = io.ReadAll(inputReader)
		if err != nil || len(validationBytes) == 0 {
			message.Fatalf(err, "error reading from stdin: %v", err)
		}
	} else if !strings.HasSuffix(path, ".yaml") {
		message.Fatalf(fmt.Errorf("input file must be a yaml file"), "input file must be a yaml file")
	} else {
		// Read the validation file
		validationBytes, err = common.ReadFileToBytes(path)
		if err != nil {
			message.Fatalf(err, "error reading file: %v", err)
		}
	}
	return validationBytes, nil
}

// RunSingleValidation runs a single validation
func RunSingleValidation(ctx context.Context, validationBytes []byte, opts ...types.LulaValidationOption) (lulaValidation types.LulaValidation, err error) {
	var validation common.Validation

	err = yaml.Unmarshal(validationBytes, &validation)
	if err != nil {
		return lulaValidation, err
	}

	lulaValidation, err = validation.ToLulaValidation("")
	if err != nil {
		return lulaValidation, err
	}

	err = lulaValidation.Validate(ctx, opts...)
	if err != nil {
		return lulaValidation, err
	}

	return lulaValidation, nil
}

// Provides basic templating wrapper for "all" render type
func DevTemplate(validationBytes []byte, setOpts []string) ([]byte, error) {
	// Get overrides from --set flag
	overrides, err := cmdCommon.ParseTemplateOverrides(setOpts)
	if err != nil {
		return nil, fmt.Errorf("error parsing template overrides: %v", err)
	}

	// Handles merging viper config file data + environment variables
	// Throws an error if config keys are invalid for templating
	templateData, err := template.CollectTemplatingData(cmdCommon.TemplateConstants, cmdCommon.TemplateVariables, overrides)
	if err != nil {
		return nil, fmt.Errorf("error collecting templating data: %v", err)
	}

	templateRenderer := template.NewTemplateRenderer(templateData)
	output, err := templateRenderer.Render(string(validationBytes), "all")
	if err != nil {
		return nil, fmt.Errorf("error rendering template: %v", err)
	}

	return output, nil
}

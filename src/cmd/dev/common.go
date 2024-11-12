package dev

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/defenseunicorns/lula/src/config"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

const STDIN = "0"
const NO_TIMEOUT = -1
const DEFAULT_TIMEOUT = 1

var devCmd = &cobra.Command{
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

type flags struct {
	InputFile        string // -f --input-file
	OutputFile       string // -o --output-file
	Timeout          int    // -t --timeout
	ConfirmExecution bool   // --confirm-execution
}

var RunInteractively bool = true // default to run dev command interactively

// Include adds the tools command to the root command.
func Include(rootCmd *cobra.Command) {
	rootCmd.AddCommand(devCmd)
}

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

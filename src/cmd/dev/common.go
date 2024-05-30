package dev

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var devCmd = &cobra.Command{
	Use:     "dev",
	Aliases: []string{"t"},
	Short:   "Collection of dev commands to make dev life easier",
}

// Include adds the tools command to the root command.
func Include(rootCmd *cobra.Command) {
	rootCmd.AddCommand(devCmd)
}

// ReadValidation reads the validation yaml file and returns the validation bytes
func ReadValidation(cmd *cobra.Command, spinner *message.Spinner, path string, timeout int) ([]byte, error) {
	var validationBytes []byte
	var err error

	if validateOpts.InputFile == STDIN {
		var inputReader io.Reader = cmd.InOrStdin()

		// If the timeout is not -1, wait for the timeout then close and return an error
		go func() {
			if validateOpts.Timeout != NO_TIMEOUT {
				time.Sleep(time.Duration(validateOpts.Timeout) * time.Second)
				cmd.Help()
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
	} else if !strings.HasSuffix(validateOpts.InputFile, ".yaml") {
		message.Fatalf(fmt.Errorf("input file must be a yaml file"), "input file must be a yaml file")
	} else {
		// Read the validation file
		validationBytes, err = common.ReadFileToBytes(validateOpts.InputFile)
		if err != nil {
			message.Fatalf(err, "error reading file: %v", err)
		}
	}
	return validationBytes, nil
}

// RunSingleValidation runs a single validation
func RunSingleValidation(validationBytes []byte, opts ...types.LulaValidationOption) (lulaValidation types.LulaValidation, err error) {
	var validation common.Validation

	err = yaml.Unmarshal(validationBytes, &validation)
	if err != nil {
		return lulaValidation, err
	}

	lulaValidation, err = validation.ToLulaValidation()
	if err != nil {
		return lulaValidation, err
	}

	err = lulaValidation.Validate(opts...)
	if err != nil {
		return lulaValidation, err
	}

	return lulaValidation, nil
}

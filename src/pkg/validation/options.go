package validation

import (
	"github.com/defenseunicorns/lula/src/pkg/message"
)

type Option func(*Validator) error

func WithAllowExecution(confirmExecution, runNonInteractively bool) Option {
	return func(v *Validator) error {
		if !confirmExecution {
			if !runNonInteractively {
				v.requestExecutionConfirmation = true
			} else {
				message.Infof("Validations requiring execution will NOT be run")
			}
		} else {
			v.runExecutableValidations = true
		}
		return nil
	}
}

func WithSaveResources(saveResources bool, outputDir string) Option {
	return func(v *Validator) error {
		if saveResources {
			v.saveResources = true
			v.outputDir = outputDir
		}

		return nil
	}
}

func WithStrict(strict bool) Option {
	return func(v *Validator) error {
		v.strict = strict
		return nil
	}
}

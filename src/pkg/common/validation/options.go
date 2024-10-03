package validation

import (
	"fmt"

	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

type Option func(*ValidationContext) error

func WithCompositionContext(cctx *composition.CompositionContext, path string) Option {
	return func(ctx *ValidationContext) error {
		var err error
		if cctx == nil {
			cctx, err = composition.New(composition.WithModelFromLocalPath(path))
			if err != nil {
				return fmt.Errorf("error creating composition context: %v", err)
			}
		}
		ctx.cctx = cctx
		return nil
	}
}

func WithAllowExecution(confirmExecution, runNonInteractively bool) Option {
	return func(ctx *ValidationContext) error {
		if !confirmExecution {
			if !runNonInteractively {
				ctx.requestExecutionConfirmation = true
			} else {
				message.Infof("Validations requiring execution will NOT be run")
			}
		} else {
			ctx.runExecutableValidations = true
		}
		return nil
	}
}

func WithResourcesDir(saveResources bool, rootDir string) Option {
	return func(ctx *ValidationContext) error {
		if saveResources {
			ctx.resourcesDir = rootDir
		}
		return nil
	}
}

package validate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/validate"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/stretchr/testify/require"
)

var (
	validInputFile    = "../../test/unit/common/oscal/valid-component.yaml"
	invalidOutputFile = "../../test/unit/common/validation/validation.opa.yaml"
)

func TestValidate(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.yaml")
	ctx := context.Background()

	t.Run("Test validate with valid input file", func(t *testing.T) {
		err := validate.Validate(ctx, validInputFile, outputFile, "")
		require.Nil(t, err)

		compiledBytes, err := os.ReadFile(outputFile)
		require.Nilf(t, err, "error reading file: %v", err)

		_, err = oscal.NewOscalModel(compiledBytes)
		require.Nilf(t, err, "error creating oscal model from validated component definition: %v", err)
	})

	t.Run("Test validate with invalid input file", func(t *testing.T) {
		err := validate.Validate(ctx, "not-a-file.yaml", "output.yaml", "")
		require.NotNil(t, err)
		require.ErrorContains(t, err, "error validating on path")
	})

	t.Run("Test validate with invalid input file and context", func(t *testing.T) {
		err := validate.Validate(ctx, "not-a-file.yaml", "output.yaml", "",
			validation.WithCompositionContext(nil, "not-a-file.yaml"))
		require.NotNil(t, err)
		require.ErrorContains(t, err, "error creating validation context")
	})

	t.Run("Test validate with invalid output file", func(t *testing.T) {
		err := validate.Validate(ctx, validInputFile, invalidOutputFile, "")
		require.NotNil(t, err)
		require.ErrorContains(t, err, "error writing component to file")
	})

	t.Run("Test validate with invalid target", func(t *testing.T) {
		err := validate.Validate(ctx, validInputFile, "output.yaml", "invalid-target")
		require.NotNil(t, err)
		require.ErrorContains(t, err, "error validating on path")
	})

}

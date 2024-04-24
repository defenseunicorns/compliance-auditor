package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd/tools"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

var (
	validInputFile   = "../../test/unit/common/compilation/component-definition-local-and-remote.yaml"
	invalidInputFile = "../../test/unit/common/valid-api-spec.yaml"
)

func TestCompileComponentDefinition(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.yaml")

	t.Run("compiles valid component definition", func(t *testing.T) {
		err := tools.Compile(validInputFile, outputFile)
		if err != nil {
			t.Fatalf("error compiling component definition: %s", err)
		}

		compiledBytes, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("error reading compiled component definition: %s", err)
		}
		compiledModel, err := oscal.NewOscalModel(compiledBytes)
		if err != nil {
			t.Fatalf("error creating oscal model from compiled component definition: %s", err)
		}

		if compiledModel.ComponentDefinition.BackMatter.Resources == nil {
			t.Fatal("compiled component definition is nil")
		}

		if len(*compiledModel.ComponentDefinition.BackMatter.Resources) <= 1 {
			t.Fatalf("expected 2 resources, got %d", len(*compiledModel.ComponentDefinition.BackMatter.Resources))
		}
	})

	t.Run("invalid component definition throws error", func(t *testing.T) {
		err := tools.Compile(invalidInputFile, outputFile)
		if err == nil {
			t.Fatal("expected error compiling invalid component definition")
		}
	})
}

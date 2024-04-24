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

func TestComposeComponentDefinition(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.yaml")

	t.Run("composes valid component definition", func(t *testing.T) {
		err := tools.Compose(validInputFile, outputFile)
		if err != nil {
			t.Fatalf("error composing component definition: %s", err)
		}

		compiledBytes, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("error reading composed component definition: %s", err)
		}
		compiledModel, err := oscal.NewOscalModel(compiledBytes)
		if err != nil {
			t.Fatalf("error creating oscal model from composed component definition: %s", err)
		}

		if compiledModel.ComponentDefinition.BackMatter.Resources == nil {
			t.Fatal("composed component definition is nil")
		}

		if len(*compiledModel.ComponentDefinition.BackMatter.Resources) <= 1 {
			t.Fatalf("expected 2 resources, got %d", len(*compiledModel.ComponentDefinition.BackMatter.Resources))
		}
	})

	t.Run("invalid component definition throws error", func(t *testing.T) {
		err := tools.Compose(invalidInputFile, outputFile)
		if err == nil {
			t.Fatal("expected error composing invalid component definition")
		}
	})
}

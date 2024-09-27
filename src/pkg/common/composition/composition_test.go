package composition_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"gopkg.in/yaml.v3"
)

const (
	allRemote          = "../../../test/e2e/scenarios/validation-composition/component-definition.yaml"
	allRemoteBadHref   = "../../../test/e2e/scenarios/validation-composition/component-definition-bad-href.yaml"
	allLocal           = "../../../test/unit/common/composition/component-definition-all-local.yaml"
	allLocalBadHref    = "../../../test/unit/common/composition/component-definition-all-local-bad-href.yaml"
	localAndRemote     = "../../../test/unit/common/composition/component-definition-local-and-remote.yaml"
	subComponentDef    = "../../../test/unit/common/composition/component-definition-import-compdefs.yaml"
	compDefMultiImport = "../../../test/unit/common/composition/component-definition-import-multi-compdef.yaml"

	// TODO: add tests for templating
	compDefNestedImport = "../../../test/unit/common/composition/component-definition-import-nested-compdef.yaml"
	compDefMultiTmpl    = "../../../test/unit/common/composition/component-definition-local-and-remote-template.yaml"

	// Also, add cmd tests...? compare golden composed file?
)

func TestComposeFromPath(t *testing.T) {
	test := func(t *testing.T, path string, opts ...composition.Option) (*oscalTypes_1_1_2.OscalCompleteSchema, error) {
		t.Helper()
		ctx := context.Background()

		options := append([]composition.Option{composition.WithModelFromLocalPath(path)}, opts...)
		cc, err := composition.New(options...)
		if err != nil {
			return nil, err
		}

		model, err := cc.ComposeFromPath(ctx, path)
		if err != nil {
			return nil, err
		}

		return model, nil
	}

	t.Run("No imports, local validations", func(t *testing.T) {
		model, err := test(t, allLocal)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("No imports, local validations, bad href", func(t *testing.T) {
		model, err := test(t, allLocalBadHref)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("No imports, remote validations", func(t *testing.T) {
		model, err := test(t, allRemote)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("Imports, no components", func(t *testing.T) {
		model, err := test(t, allRemote)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("No imports, bad remote validations", func(t *testing.T) {
		model, err := test(t, allRemoteBadHref)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("Errors when file does not exist", func(t *testing.T) {
		_, err := test(t, "nonexistent")
		if err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("Resolves relative paths", func(t *testing.T) {
		model, err := test(t, localAndRemote)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})
}

func TestComposeComponentDefinitions(t *testing.T) {
	test := func(t *testing.T, compDef *oscalTypes_1_1_2.ComponentDefinition, path string, opts ...composition.Option) (*oscalTypes_1_1_2.OscalCompleteSchema, error) {
		t.Helper()
		ctx := context.Background()

		options := append([]composition.Option{composition.WithModelFromLocalPath(path)}, opts...)
		cc, err := composition.New(options...)
		if err != nil {
			return nil, err
		}

		baseDir := filepath.Dir(path)

		err = cc.ComposeComponentDefinitions(ctx, compDef, baseDir)
		if err != nil {
			return nil, err
		}

		return &oscalTypes_1_1_2.OscalCompleteSchema{
			ComponentDefinition: compDef,
		}, nil
	}

	t.Run("No imports, local validations", func(t *testing.T) {
		og := getComponentDef(allLocal, t)
		compDef := getComponentDef(allLocal, t)

		model, err := test(t, compDef, allLocal)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		if compDefComposed == nil {
			t.Error("expected the component definition to be non-nil")
		}

		// Only the last-modified timestamp should be different
		if !reflect.DeepEqual(*og.BackMatter, *compDefComposed.BackMatter) {
			t.Error("expected the back matter to be unchanged")
		}
	})

	t.Run("No imports, remote validations", func(t *testing.T) {
		og := getComponentDef(allRemote, t)
		compDef := getComponentDef(allRemote, t)

		model, err := test(t, compDef, allRemote)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		if compDefComposed == nil {
			t.Error("expected the component definition to be non-nil")
		}

		if reflect.DeepEqual(*og, *compDefComposed) {
			t.Errorf("expected component definition to have changed.")
		}
	})

	t.Run("Imports, no components", func(t *testing.T) {
		og := getComponentDef(subComponentDef, t)
		compDef := getComponentDef(subComponentDef, t)

		model, err := test(t, compDef, subComponentDef)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		if compDefComposed == nil {
			t.Error("expected the component definition to be non-nil")
		}

		if compDefComposed.Components == og.Components {
			t.Error("expected there to be components")
		}

		if compDefComposed.BackMatter == og.BackMatter {
			t.Error("expected the back matter to be changed")
		}
	})

	t.Run("imports, no components, multiple component definitions from import", func(t *testing.T) {
		og := getComponentDef(compDefMultiImport, t)
		compDef := getComponentDef(compDefMultiImport, t)

		model, err := test(t, compDef, subComponentDef)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		if compDefComposed == nil {
			t.Error("expected the component definition to be non-nil")
		}

		if compDefComposed.Components == og.Components {
			t.Error("expected there to be components")
		}

		if compDefComposed.BackMatter == og.BackMatter {
			t.Error("expected the back matter to be changed")
		}

		if len(*compDefComposed.Components) != 1 {
			t.Error("expected there to be 1 component")
		}
	})

}

func TestCompileComponentValidations(t *testing.T) {
	test := func(t *testing.T, compDef *oscalTypes_1_1_2.ComponentDefinition, path string, opts ...composition.Option) (*oscalTypes_1_1_2.OscalCompleteSchema, error) {
		t.Helper()
		ctx := context.Background()

		options := append([]composition.Option{composition.WithModelFromLocalPath(path)}, opts...)
		cc, err := composition.New(options...)
		if err != nil {
			return nil, err
		}

		baseDir := filepath.Dir(path)

		err = cc.ComposeComponentValidations(ctx, compDef, baseDir)
		if err != nil {
			return nil, err
		}

		return &oscalTypes_1_1_2.OscalCompleteSchema{
			ComponentDefinition: compDef,
		}, nil
	}

	t.Run("all local", func(t *testing.T) {
		og := getComponentDef(allLocal, t)
		compDef := getComponentDef(allLocal, t)

		model, err := test(t, compDef, allLocal)
		if err != nil {
			t.Fatalf("error composing validations: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		if compDefComposed == nil {
			t.Error("expected the component definition to be non-nil")
		}

		// Only the last-modified timestamp should be different
		if !reflect.DeepEqual(*og.BackMatter, *compDefComposed.BackMatter) {
			t.Error("expected the back matter to be unchanged")
		}
	})

	t.Run("all remote", func(t *testing.T) {
		og := getComponentDef(allRemote, t)
		compDef := getComponentDef(allRemote, t)

		model, err := test(t, compDef, allRemote)
		if err != nil {
			t.Fatalf("error composing validations: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		if compDefComposed == nil {
			t.Error("expected the component definition to be non-nil")
		}

		if reflect.DeepEqual(*og, *compDefComposed) {
			t.Error("expected the component definition to be changed")
		}

		if compDefComposed.BackMatter == nil {
			t.Error("expected the component definition to have back matter")
		}

		if og.Metadata.LastModified == compDefComposed.Metadata.LastModified {
			t.Error("expected the component definition to have a different last modified timestamp")
		}
	})

	t.Run("local and remote", func(t *testing.T) {
		og := getComponentDef(localAndRemote, t)
		compDef := getComponentDef(localAndRemote, t)

		model, err := test(t, compDef, localAndRemote)
		if err != nil {
			t.Fatalf("error composing validations: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		if compDefComposed == nil {
			t.Error("expected the component definition to be non-nil")
		}

		if reflect.DeepEqual(*og, *compDefComposed) {
			t.Error("expected the component definition to be changed")
		}
	})
}

func getComponentDef(path string, t *testing.T) *oscalTypes_1_1_2.ComponentDefinition {
	compDef, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error reading component definition file: %v", err)
	}

	var oscalModel oscalTypes_1_1_2.OscalModels
	if err := yaml.Unmarshal(compDef, &oscalModel); err != nil {
		t.Fatalf("Error unmarshalling component definition: %v", err)
	}
	return oscalModel.ComponentDefinition
}

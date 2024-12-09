package composition_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	allRemote           = "../../../test/e2e/scenarios/validation-composition/component-definition.yaml"
	allRemoteBadHref    = "../../../test/e2e/scenarios/validation-composition/component-definition-bad-href.yaml"
	allLocal            = "../../../test/unit/common/composition/component-definition-all-local.yaml"
	allLocalBadHref     = "../../../test/unit/common/composition/component-definition-all-local-bad-href.yaml"
	localAndRemote      = "../../../test/unit/common/composition/component-definition-local-and-remote.yaml"
	subComponentDef     = "../../../test/unit/common/composition/component-definition-import-compdefs.yaml"
	compDefMultiImport  = "../../../test/unit/common/composition/component-definition-import-multi-compdef.yaml"
	compDefNestedImport = "../../../test/unit/common/composition/component-definition-import-nested-compdef.yaml"
	compDefTmpl         = "../../../test/unit/common/composition/component-definition-template.yaml"
	compDefNestedTmpl   = "../../../test/unit/common/composition/component-definition-import-nested-compdef-template.yaml"
)

func TestComposeFromPath(t *testing.T) {
	test := func(t *testing.T, path string, opts ...composition.Option) (*oscalTypes.OscalCompleteSchema, error) {
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

	t.Run("Nested imports, no components", func(t *testing.T) {
		model, err := test(t, compDefNestedImport)
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

	t.Run("Templated component definition, error", func(t *testing.T) {
		model, err := test(t, compDefTmpl)
		if err == nil {
			t.Fatalf("Should encounter error composing component definitions: %v", err)
		}
		if model != nil {
			t.Error("expected the model not to be composed")
		}
	})

	// Test the templating of the component definition where the validation is not rendered -> empty resources in backmatter
	t.Run("Templated component definition with nested imports, validations not rendered - no resources", func(t *testing.T) {
		tmplOpts := []composition.Option{
			composition.WithRenderSettings("constants", true),
			composition.WithTemplateRenderer("constants", map[string]interface{}{
				"templated_comp_def": interface{}("component-definition-template.yaml"),
				"type":               interface{}("software"),
				"title":              interface{}("lula"),
			}, []template.VariableConfig{}, []string{}),
		}

		model, err := test(t, compDefTmpl, tmplOpts...)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		require.NotNil(t, compDefComposed)
		require.NotNil(t, compDefComposed.Components)
		require.NotNil(t, compDefComposed.BackMatter)
		require.NotNil(t, compDefComposed.BackMatter.Resources)
		require.Equal(t, len(*compDefComposed.BackMatter.Resources), 0)
	})

	// Test the templating of the component definition with nested templated imports
	t.Run("Templated component definition with nested imports, validations rendered", func(t *testing.T) {
		tmplOpts := []composition.Option{
			composition.WithRenderSettings("constants", true),
			composition.WithTemplateRenderer("constants", map[string]interface{}{
				"templated_comp_def": interface{}("component-definition-template.yaml"),
				"type":               interface{}("software"),
				"title":              interface{}("lula"),
				"resources": interface{}(map[string]interface{}{
					"name":      interface{}("test-pod-label"),
					"namespace": interface{}("validation-test"),
					"exemptions": []interface{}{
						interface{}("one"),
						interface{}("two"),
						interface{}("three"),
					},
				}),
			}, []template.VariableConfig{}, []string{}),
		}
		model, err := test(t, compDefNestedTmpl, tmplOpts...)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		require.NotNil(t, compDefComposed)
		require.NotNil(t, compDefComposed.Components)
		require.NotNil(t, compDefComposed.BackMatter)
		require.NotNil(t, compDefComposed.BackMatter.Resources)
		require.Len(t, *compDefComposed.BackMatter.Resources, 1)
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
	test := func(t *testing.T, compDef *oscalTypes.ComponentDefinition, path string, opts ...composition.Option) (*oscalTypes.OscalCompleteSchema, error) {
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

		return &oscalTypes.OscalCompleteSchema{
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
		require.NotNil(t, compDefComposed)

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
		require.NotNil(t, compDefComposed)

		if reflect.DeepEqual(*og, *compDefComposed) {
			t.Error("expected component definition to have changed.")
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
		require.NotNil(t, compDefComposed)

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

		model, err := test(t, compDef, compDefMultiImport)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		require.NotNil(t, compDefComposed)

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

	// Both "imported" components have the same component (by UUID), so those are merged
	// Both components have the same control-impementation (by control ID, not UUID), those are merged
	// All validations are linked to that single control-implementation
	t.Run("nested imports, directory changes", func(t *testing.T) {
		og := getComponentDef(compDefNestedImport, t)
		compDef := getComponentDef(compDefNestedImport, t)

		model, err := test(t, compDef, compDefNestedImport)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		compDefComposed := model.ComponentDefinition
		require.NotNil(t, compDefComposed)

		if compDefComposed.Components == og.Components {
			t.Error("expected there to be new components")
		}

		if compDefComposed.BackMatter == og.BackMatter {
			t.Error("expected the back matter to be changed")
		}

		components := *compDefComposed.Components
		if len(components) != 1 {
			t.Error("expected there to be 1 component")
		}

		if len(*components[0].ControlImplementations) != 1 {
			t.Error("expected there to be 1 control implementation")
		}

		if len(*compDefComposed.BackMatter.Resources) != 7 {
			t.Error("expected the back matter to contain 7 resources (validations)")
		}
	})

}

func TestComposeComponentValidations(t *testing.T) {
	test := func(t *testing.T, compDef *oscalTypes.ComponentDefinition, path string, opts ...composition.Option) (*oscalTypes.OscalCompleteSchema, error) {
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

		return &oscalTypes.OscalCompleteSchema{
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
		require.NotNil(t, compDefComposed)

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
		require.NotNil(t, compDefComposed)

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
		require.NotNil(t, compDefComposed)

		if reflect.DeepEqual(*og, *compDefComposed) {
			t.Error("expected the component definition to be changed")
		}
	})
}

func getComponentDef(path string, t *testing.T) *oscalTypes.ComponentDefinition {
	compDef, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error reading component definition file: %v", err)
	}

	var oscalModel oscalTypes.OscalModels
	if err := yaml.Unmarshal(compDef, &oscalModel); err != nil {
		t.Fatalf("Error unmarshalling component definition: %v", err)
	}
	return oscalModel.ComponentDefinition
}

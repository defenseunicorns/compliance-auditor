package composition_test

import (
	"os"
	"reflect"
	"testing"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
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
)

func TestComposeFromPath(t *testing.T) {
	t.Run("No imports, local validations", func(t *testing.T) {
		model, err := composition.ComposeFromPath(allLocal)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("No imports, local validations, bad href", func(t *testing.T) {
		model, err := composition.ComposeFromPath(allLocalBadHref)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("No imports, remote validations", func(t *testing.T) {
		model, err := composition.ComposeFromPath(allRemote)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("No imports, bad remote validations", func(t *testing.T) {
		model, err := composition.ComposeFromPath(allRemoteBadHref)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})

	t.Run("Errors when file does not exist", func(t *testing.T) {
		_, err := composition.ComposeFromPath("nonexistent")
		if err == nil {
			t.Error("expected an error")
		}
	})

	t.Run("Resolves relative paths", func(t *testing.T) {
		model, err := composition.ComposeFromPath(localAndRemote)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if model == nil {
			t.Error("expected the model to be composed")
		}
	})
}

func TestComposeComponentDefinitions(t *testing.T) {
	t.Run("No imports, local validations", func(t *testing.T) {
		og := getComponentDef(allLocal, t)
		compDef := getComponentDef(allLocal, t)
		reset, err := common.SetCwdToFileDir(allLocal)
		defer reset()
		if err != nil {
			t.Fatalf("Error setting cwd to file dir: %v", err)
		}
		err = composition.ComposeComponentDefinitions(compDef)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		// Only the last-modified timestamp should be different
		if !reflect.DeepEqual(*og.BackMatter, *compDef.BackMatter) {
			t.Error("expected the back matter to be unchanged")
		}
	})

	t.Run("No imports, remote validations", func(t *testing.T) {
		og := getComponentDef(allRemote, t)
		compDef := getComponentDef(allRemote, t)
		reset, err := common.SetCwdToFileDir(allRemote)
		defer reset()
		if err != nil {
			t.Fatalf("Error setting cwd to file dir: %v", err)
		}
		err = composition.ComposeComponentDefinitions(compDef)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		if reflect.DeepEqual(*og, *compDef) {
			t.Errorf("expected component definition to have changed.")
		}
	})

	t.Run("Imports, no components", func(t *testing.T) {
		og := getComponentDef(subComponentDef, t)
		compDef := getComponentDef(subComponentDef, t)
		reset, err := common.SetCwdToFileDir(subComponentDef)
		defer reset()
		if err != nil {
			t.Fatalf("Error setting cwd to file dir: %v", err)
		}
		err = composition.ComposeComponentDefinitions(compDef)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}

		if compDef.Components == og.Components {
			t.Error("expected there to be components")
		}

		if compDef.BackMatter == og.BackMatter {
			t.Error("expected the back matter to be changed")
		}
	})

	t.Run("imports, no components, multiple component definitions from import", func(t *testing.T) {
		og := getComponentDef(compDefMultiImport, t)
		compDef := getComponentDef(compDefMultiImport, t)
		reset, err := common.SetCwdToFileDir(compDefMultiImport)
		defer reset()
		if err != nil {
			t.Fatalf("Error setting cwd to file dir: %v", err)
		}
		err = composition.ComposeComponentDefinitions(compDef)
		if err != nil {
			t.Fatalf("Error composing component definitions: %v", err)
		}
		if compDef.Components == og.Components {
			t.Error("expected there to be components")
		}

		if compDef.BackMatter == og.BackMatter {
			t.Error("expected the back matter to be changed")
		}

		if len(*compDef.Components) != 1 {
			t.Error("expected there to be 2 components")
		}
	})

}

func TestCompileComponentValidations(t *testing.T) {

	t.Run("all local", func(t *testing.T) {
		og := getComponentDef(allLocal, t)
		compDef := getComponentDef(allLocal, t)
		reset, err := common.SetCwdToFileDir(allLocal)
		defer reset()
		if err != nil {
			t.Fatalf("Error setting cwd to file dir: %v", err)
		}
		err = composition.ComposeComponentValidations(compDef)
		if err != nil {
			t.Fatalf("Error compiling component validations: %v", err)
		}

		// Only the last-modified timestamp should be different
		if !reflect.DeepEqual(*og.BackMatter, *compDef.BackMatter) {
			t.Error("expected the back matter to be unchanged")
		}
	})

	t.Run("all remote", func(t *testing.T) {
		og := getComponentDef(allRemote, t)
		compDef := getComponentDef(allRemote, t)
		reset, err := common.SetCwdToFileDir(allRemote)
		defer reset()
		if err != nil {
			t.Fatalf("Error setting cwd to file dir: %v", err)
		}
		err = composition.ComposeComponentValidations(compDef)
		if err != nil {
			t.Fatalf("Error compiling component validations: %v", err)
		}
		if reflect.DeepEqual(*og, *compDef) {
			t.Error("expected the component definition to be changed")
		}

		if compDef.BackMatter == nil {
			t.Error("expected the component definition to have back matter")
		}

		if og.Metadata.LastModified == compDef.Metadata.LastModified {
			t.Error("expected the component definition to have a different last modified timestamp")
		}
	})

	t.Run("local and remote", func(t *testing.T) {
		og := getComponentDef(localAndRemote, t)
		compDef := getComponentDef(localAndRemote, t)
		reset, err := common.SetCwdToFileDir(localAndRemote)
		defer reset()
		if err != nil {
			t.Fatalf("Error setting cwd to file dir: %v", err)
		}
		err = composition.ComposeComponentValidations(compDef)
		if err != nil {
			t.Fatalf("Error compiling component validations: %v", err)
		}

		if reflect.DeepEqual(*og, *compDef) {
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

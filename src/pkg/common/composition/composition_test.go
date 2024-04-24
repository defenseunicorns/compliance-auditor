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
	allRemote      = "../../../test/e2e/scenarios/validation-composition/component-definition.yaml"
	allLocal       = "../../../test/unit/common/compilation/component-definition-all-local.yaml"
	localAndRemote = "../../../test/unit/common/compilation/component-definition-local-and-remote.yaml"
)

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

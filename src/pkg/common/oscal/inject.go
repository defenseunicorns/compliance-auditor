package oscal

import (
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// InjectJSONPathValues injects values into an OSCAL model using JSONPath
func InjectJSONPathValues(model map[string]interface{}, path string, values map[string]interface{}) error {
	// Right node is the model
	// Left node is an empty oscalTypes_1_1_2.OscalModels with the values injected at the path specified
	// Merge the two nodes, that becomes model

	modelNode, err := yaml.FromMap(model)
	if err != nil {
		return fmt.Errorf("failed to create left node from values: %v", err)
	}

	valuesNode, err := yaml.FromMap(values)
	if err != nil {
		return fmt.Errorf("failed to create left node from values: %v", err)
	}
	fmt.Printf("valuesNode kind: %v\n", valuesNode.GetKind())

	injectionNode, err := modelNode.Pipe(yaml.Lookup(splitPath(path)...))
	if err != nil {
		return fmt.Errorf("failed to find path %s in model: %v", path, err)
	}
	fmt.Printf("injectionNode kind: %v\n", injectionNode.GetKind())

	// merge3.Merge(injectionNode, valuesNode, mergedNode)

	// merge lNode into injectionNode -> put injectionNode back into rNode?

	return nil
}

// splitPath splits a path by '.' into a path array - is there a lib function to do this and possibly handle things like [] or escaped '.'
func splitPath(path string) []string {
	return strings.Split(path, ".")
}

package transform

import (
	"fmt"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/utils"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// InjectMapData injects the subset map into a target map at the path
// TODO: should this behave differently if the path is not found? Or if you want to replace a seq instead of append?
func InjectMapData(target, subset map[string]interface{}, path string) (map[string]interface{}, error) {
	pathSlice := utils.SmarterPathSplitter(path, ".")

	// Convert the target and subset maps to yaml nodes
	targetNode, err := yaml.FromMap(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create target node from map: %v", err)
	}

	subsetNode, err := yaml.FromMap(subset)
	if err != nil {
		return nil, fmt.Errorf("failed to create subset node from map: %v", err)
	}

	// Build the target filters, update the target with subset data
	filters, err := BuildFilters(targetNode, pathSlice)
	if err != nil {
		return nil, err
	}

	targetSubsetNode, err := targetNode.Pipe(filters...)
	if err != nil {
		return nil, fmt.Errorf("error identifying subset node: %v", err)
	}

	// Alternate merge based on custom merge function
	// TODO: add option to replace all and use the kyaml merge function?
	err = mergeYAMLNodes(targetSubsetNode, subsetNode)
	if err != nil {
		return nil, fmt.Errorf("error merging subset into target: %v", err)
	}

	// Inject the updated node back into targetNode
	if len(pathSlice) == 0 {
		targetNode = targetSubsetNode
	} else {
		if err = SetNodeAtPath(targetNode, targetSubsetNode, filters, pathSlice); err != nil {
			return nil, fmt.Errorf("error setting merged node back into target: %v", err)
		}
	}

	// Write targetNode into map[string]interface{}
	var targetMap map[string]interface{}
	targetNode.YNode().Decode(&targetMap)

	return targetMap, nil
}

// TODO: add support to delete a field
// func EjectMapData(target map[string]interface{}, path string) (map[string]interface{}, error) {
// 	pathSlice, err := splitPath(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to split path: %v", err)
// 	}

// 	// Create a new node from the target map
// 	targetNode, err := yaml.FromMap(target)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create target node from map: %v", err)
// 	}

// 	filters, err := buildFilters(targetNode, pathSlice)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// filters = append(filters, yaml.FieldClearer{Name: pathSlice[len(pathSlice)-1]})

// 	targetSubsetNode, err := targetNode.Pipe(filters...)
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding subset node in target: %v", err)
// 	}

// 	// merge it back in?

// 	var targetMap map[string]interface{}
// 	targetNode.YNode().Decode(&targetMap)

// 	return targetMap, nil
// }

// TODO: add support to add an element or map
// func ModifyMapValue(target map[string]interface{}, path string, value interface{}) (map[string]interface{}, error) {
// 	pathSlice, err := splitPath(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to split path: %v", err)
// 	}

// 	// Create a new node from the target map
// 	targetNode, err := yaml.FromMap(target)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create target node from map: %v", err)
// 	}

// 	filters, err := buildFilters(pathSlice)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// add a elementsetter?
// }

// ExtractMapData extracts a subset of data from a map[string]interface{} and returns it as a map[string]interface{}
func ExtractMapData(target map[string]interface{}, path string) (map[string]interface{}, error) {
	pathSlice, err := splitPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to split path: %v", err)
	}

	// Convert the target and subset maps to yaml nodes
	targetNode, err := yaml.FromMap(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create target node from map: %v", err)
	}

	// Build the target filters, update the target with subset data
	filters, err := BuildFilters(targetNode, pathSlice)
	if err != nil {
		return nil, err
	}
	targetSubsetNode, err := targetNode.Pipe(filters...)
	if err != nil {
		return nil, fmt.Errorf("error identifying subset node: %v", err)
	}

	// Write targetSubsetNode into map[string]interface{}
	targetSubsetMap, err := targetSubsetNode.Map()
	if err != nil {
		return nil, fmt.Errorf("failed to convert target subset node to map: %v", err)
	}

	return targetSubsetMap, nil
}

// splitPath splits a path by '.' into a path array and handles filters in the path
func splitPath(path string) ([]string, error) {
	if path == "" {
		return []string{}, nil
	}

	// Regex to match path segments including filters
	re := regexp.MustCompile(`[^.\[\]]+(?:\[[^\[\]]+\])?`)
	matches := re.FindAllString(path, -1)
	if matches == nil {
		return nil, fmt.Errorf("invalid path format")
	}
	return matches, nil
}

// stringToMap converts a dot notation string like "metadata.name=foo" into a map[string]interface{}
func stringToMap(input string) map[string]interface{} {
	// Split the input string into the key part and the value part
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return nil
	}

	keyPart := parts[0]
	valuePart := parts[1]

	// Split the key part on "." to handle nested maps
	keys := strings.Split(keyPart, ".")

	// Create a nested map structure based on the keys
	result := make(map[string]interface{})
	current := result

	for i, key := range keys {
		if i == len(keys)-1 {
			// Last key, set the value
			current[key] = valuePart
		} else {
			// If the key doesn't exist yet, create a new nested map
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]interface{})
			}
			// Move to the next level
			current = current[key].(map[string]interface{})
		}
	}

	return result
}

package inject

import (
	"fmt"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/utils"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type filterParts struct {
	key   string
	value string
}

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
	filters, err := buildFilters(targetNode, pathSlice)
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
		if err = setNodeAtPath(targetNode, targetSubsetNode, filters, pathSlice); err != nil {
			return nil, fmt.Errorf("error setting merged node back into target: %v", err)
		}
	}

	// Write targetNode into map[string]interface{}
	var targetMap map[string]interface{}
	targetNode.YNode().Decode(&targetMap)

	return targetMap, nil
}

func GetRNode(target map[string]interface{}, path string) (*yaml.RNode, error) {
	pathSlice := utils.SmarterPathSplitter(path, ".")

	// Convert the target to yaml nodes
	targetNode, err := yaml.FromMap(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create target node from map: %v", err)
	}

	// Build the target filters, update the target with subset data
	filters, err := buildFilters(targetNode, pathSlice)
	if err != nil {
		return nil, err
	}

	targetSubsetNode, err := targetNode.Pipe(filters...)
	if err != nil {
		return nil, fmt.Errorf("error finding subset node in target: %v", err)
	}

	return targetSubsetNode, nil
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
	filters, err := buildFilters(targetNode, pathSlice)
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

// setNodeAtPath injects the updated node into targetNode according to the specified path
func setNodeAtPath(targetNode *yaml.RNode, targetSubsetNode *yaml.RNode, filters []yaml.Filter, pathSlice []string) error {
	// Check if the last segment is a filter, changes the behavior of the set function
	lastSegment := pathSlice[len(pathSlice)-1]

	if isFilter, filterParts, err := extractFilter(pathSlice[len(pathSlice)-1]); err != nil {
		return err
	} else if isFilter {
		keys := make([]string, 0)
		values := make([]string, 0)
		for _, part := range filterParts {
			if isComposite(lastSegment) {
				// idk how to handle this... should there be a composite filter here anyway?
				return fmt.Errorf("composite filters not supported in final path segment")
			} else {
				keys = append(keys, part.key)
				values = append(values, part.value)
			}
		}
		filters = append(filters[:len(filters)-1], yaml.ElementSetter{
			Element: targetSubsetNode.Document(),
			Keys:    keys,
			Values:  values,
		})
	} else {
		filters = append(filters[:len(filters)-1], yaml.SetField(lastSegment, targetSubsetNode))
	}

	return targetNode.PipeE(filters...)
}

// returnIndexFromComplexFilters returns the index of the node that matches the filterParts
// e.g., [key1=value1,key2=value2], [composite.key=value], [val.key.test=bar]
func returnIndexFromComplexFilters(targetNode *yaml.RNode, parentFilters []yaml.Filter, filterParts []filterParts) (int, error) {
	index := -1

	parentNode, err := targetNode.Pipe(parentFilters...)
	if err != nil {
		return index, err
	}

	if parentNode.YNode().Kind == yaml.SequenceNode {
		nodes, err := parentNode.Elements()
		if err != nil {
			return index, err
		}
		for i, node := range nodes {
			if nodeMatchesAllFilters(node, filterParts) {
				index = i
				break
			}
		}
	} else {
		return index, fmt.Errorf("expected sequence node, but got %v", parentNode.YNode().Kind)
	}

	return index, nil
}

func nodeMatchesAllFilters(node *yaml.RNode, filterParts []filterParts) bool {
	for _, part := range filterParts {
		if isComposite(part.key) {
			compositeFilters := buildCompositeFilters(part.key, part.value)
			n, err := node.Pipe(compositeFilters...)
			if err != nil || n == nil {
				return false
			}
		} else {
			n, err := node.Pipe(yaml.MatchElement(part.key, part.value))
			if err != nil || n == nil {
				return false
			}
		}
	}
	return true
}

func buildFilters(targetNode *yaml.RNode, pathSlice []string) ([]yaml.Filter, error) {
	filters := make([]yaml.Filter, 0)
	for _, segment := range pathSlice {
		if isFilter, filterParts, err := extractFilter(segment); err != nil {
			return nil, err
		} else if isFilter {
			// if it's a complex filter, e.g., [key1=value1,key2=value2] or [composite.key=value], lookup the index
			if len(filterParts) > 1 || isComposite(filterParts[0].key) {
				index, err := returnIndexFromComplexFilters(targetNode, filters, filterParts)
				if err != nil {
					return nil, err
				}

				if index == -1 {
					return nil, fmt.Errorf("composite path not found")
				} else {
					filters = append(filters, yaml.GetElementByIndex(index))
				}
			} else {
				filters = append(filters, yaml.MatchElement(filterParts[0].key, filterParts[0].value))
			}
		} else {
			filters = append(filters, yaml.Lookup(segment))
		}
	}
	return filters, nil
}

func buildCompositeFilters(key, value string) []yaml.Filter {
	path := strings.Split(key, ".")
	compositeFilters := make([]yaml.Filter, 0, len(path))
	if len(path) > 1 {
		for i := 0; i < (len(path) - 1); i++ {
			compositeFilters = append(compositeFilters, yaml.Get(path[i]))
		}
	}

	compositeFilters = append(compositeFilters, yaml.MatchField(path[len(path)-1], value))
	return compositeFilters
}

// extractFilter extracts the filter parts from a string
// e.g., [key1=value1,key2=value2], [composite.key=value], [val.key.test=bar]
func extractFilter(item string) (bool, []filterParts, error) {
	if !isFilter(item) {
		return false, []filterParts{}, nil
	}
	item = strings.TrimPrefix(item, "[")
	item = strings.TrimSuffix(item, "]")

	items := strings.Split(item, ",")
	if len(items) == 0 {
		return false, []filterParts{}, fmt.Errorf("filter is empty")
	}

	filterPartsSlice := make([]filterParts, 0, len(items))
	for _, i := range items {
		if !strings.Contains(i, "=") {
			return false, []filterParts{}, fmt.Errorf("filter is not in the correct format")
		}
		filterPartsSlice = append(filterPartsSlice, filterParts{
			key:   strings.SplitN(i, "=", 2)[0],
			value: strings.SplitN(i, "=", 2)[1],
		})
	}

	return true, filterPartsSlice, nil
}

func isFilter(item string) bool {
	// check if first and last char are [ and ]
	return strings.HasPrefix(item, "[") && strings.HasSuffix(item, "]")
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

// isComposite checks if a string is a composite string, e.g., metadata.name
func isComposite(input string) bool {
	keys := strings.Split(input, ".")
	return len(keys) > 1
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

// mergeYAMLNodes recursively merges the subset node into the target node
// Note - this is an alternate to kyaml merge2 function which doesn't append lists, it replaces them
func mergeYAMLNodes(target, subset *yaml.RNode) error {
	switch subset.YNode().Kind {
	case yaml.MappingNode:
		subsetFields, err := subset.Fields()
		if err != nil {
			return err
		}
		for _, field := range subsetFields {
			subsetFieldNode, err := subset.Pipe(yaml.Lookup(field))
			if err != nil {
				return err
			}
			targetFieldNode, err := target.Pipe(yaml.Lookup(field))
			if err != nil {
				return err
			}

			if targetFieldNode == nil {
				// Field doesn't exist in target, so set it
				err = target.PipeE(yaml.SetField(field, subsetFieldNode))
				if err != nil {
					return err
				}
			} else {
				// Field exists, merge it recursively
				err = mergeYAMLNodes(targetFieldNode, subsetFieldNode)
				if err != nil {
					return err
				}
			}
		}
	case yaml.SequenceNode:
		subsetItems, err := subset.Elements()
		if err != nil {
			return err
		}
		for _, item := range subsetItems {
			target.YNode().Content = append(target.YNode().Content, item.YNode())
		}
	default:
		// Simple replacement for scalar and other nodes
		target.YNode().Value = subset.YNode().Value
	}
	return nil
}

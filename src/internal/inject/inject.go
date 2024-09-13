package inject

import (
	"fmt"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type filterParts struct {
	key   string
	value string
}

// InjectMapData injects the subset map into a target map at the path
// TODO: should this behave differently if the path is not found? Or if you want to replace a seq instead of append?
func InjectMapData(target, subset map[string]interface{}, path string) (map[string]interface{}, error) {
	pathSlice, err := splitPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to split path: %v", err)
	}

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
	filters, err := buildFilters(pathSlice)
	if err != nil {
		return nil, err
	}

	// if subset is empty, add a field clearer to the filters:
	// actually this might be better in a separate function...
	// filters = append(filters, yaml.FieldClearer{Name: pathSlice[len(pathSlice)-1]})

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
	targetMap, err := targetNode.Map()
	if err != nil {
		return nil, fmt.Errorf("failed to convert target node to map: %v", err)
	}

	return targetMap, nil
}

func ModifyMapValue(target map[string]interface{}, path string, value interface{}) (map[string]interface{}, error) {
	pathSlice, err := splitPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to split path: %v", err)
	}

	// Create a new node from the target map
	targetNode, err := yaml.FromMap(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create target node from map: %v", err)
	}

	filters, err := buildFilters(pathSlice)
	if err != nil {
		return nil, err
	}
	// add a elementsetter?
}

// DeleteMapData
// add a fieldclearer to the filters... or do above?
// func DeleteMapData(target map[string]interface{}, path string) (map[string]interface{}, error) {
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
	filters, err := buildFilters(pathSlice)
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

	if isFilter, _, filterParts, err := extractFilter(pathSlice[len(pathSlice)-1]); err != nil {
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

func buildFilters(pathSlice []string) ([]yaml.Filter, error) {
	filters := make([]yaml.Filter, 0)
	for _, segment := range pathSlice {
		if isFilter, fieldName, filterParts, err := extractFilter(segment); err != nil {
			return nil, err
		} else if isFilter {
			filters = append(filters, yaml.Lookup(fieldName))

			// Add the filters for each index selection schema
			for _, part := range filterParts {
				if isComposite(fieldName) {
					// add the fieldMatcherFilter
					valueRnode, err := yaml.FromMap(stringToMap(fmt.Sprintf("%s=%s", part.key, part.value)))
					if err != nil {
						return nil, err
					}
					// I think this field matcher only works for maps, not sequences...
					filters = append(filters, yaml.FieldMatcher{
						Name:  fieldName,
						Value: valueRnode,
					})
				} else {
					filters = append(filters, yaml.MatchElement(part.key, part.value))
				}
			}
		} else {
			filters = append(filters, yaml.Lookup(segment))
		}
	}
	return filters, nil
}

// extractFilter extracts the filter parts from a string
// e.g., fieldName[key1=value1,key2=value2], fieldName[composite.key=value]
func extractFilter(item string) (bool, string, []filterParts, error) {
	if !isFilter(item) {
		return false, "", []filterParts{}, nil
	}
	segmentParts := strings.SplitN(item, "[", 2)
	fieldName := segmentParts[0]
	filter := strings.TrimSuffix(segmentParts[1], "]")

	items := strings.Split(filter, ",")
	filterPartsSlice := make([]filterParts, 0, len(items))
	for _, item := range items {
		filterPartsSlice = append(filterPartsSlice, filterParts{
			key:   strings.SplitN(item, "=", 2)[0],
			value: strings.SplitN(item, "=", 2)[1],
		})
	}

	return true, fieldName, filterPartsSlice, nil
}

func isFilter(item string) bool {
	return strings.Contains(item, "[") && strings.Contains(item, "]")
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

// isComposite checks if a string is a composite string
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

package inject

import (
	"fmt"
	"regexp"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

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

// setNodeAtPath injects the updated node into targetNode according to the specified path
func setNodeAtPath(targetNode *yaml.RNode, targetSubsetNode *yaml.RNode, filters []yaml.Filter, pathSlice []string) error {
	// Check if the last segment is a filter, changes the behavior of the set function
	lastSegment := pathSlice[len(pathSlice)-1]

	if isFilter, _, filterParts, err := extractFilter(pathSlice[len(pathSlice)-1]); err != nil {
		return err
	} else if isFilter {
		filters = append(filters[:len(filters)-1], yaml.ElementSetter{
			Element: targetSubsetNode.Document(),
			Keys:    []string{filterParts[0]},
			Values:  []string{filterParts[1]},
		})
	} else {
		filters = append(filters[:len(filters)-1], yaml.SetField(lastSegment, targetSubsetNode))
	}

	return targetNode.PipeE(filters...)
}

func buildFilters(pathSlice []string) ([]yaml.Filter, error) {
	filters := make([]yaml.Filter, 0, len(pathSlice))
	for _, segment := range pathSlice {
		if isFilter, fieldName, filterParts, err := extractFilter(segment); err != nil {
			return nil, err
		} else if isFilter {
			filters = append(filters, yaml.Lookup(fieldName))
			filters = append(filters, yaml.MatchElement(filterParts[0], filterParts[1]))
		} else {
			filters = append(filters, yaml.Lookup(segment))
		}
	}
	return filters, nil
}

func extractFilter(item string) (bool, string, []string, error) {
	if !isFilter(item) {
		return false, "", []string{}, nil
	}
	segmentParts := strings.SplitN(item, "[", 2)
	fieldName := segmentParts[0]
	filter := strings.TrimSuffix(segmentParts[1], "]")

	filterParts := strings.SplitN(filter, "=", 2)
	if len(filterParts) != 2 {
		return false, "", []string{}, fmt.Errorf("invalid filter format: %s", item)
	}

	return true, fieldName, filterParts, nil
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
	re := regexp.MustCompile(`(?:[^.\[\]]+|\[[^\[\]]+\])+`)
	matches := re.FindAllString(path, -1)
	if matches == nil {
		return nil, fmt.Errorf("invalid path format")
	}
	return matches, nil
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

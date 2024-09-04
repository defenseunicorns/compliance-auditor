package inject

import (
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// InjectMapData injects the subset map into a target map at the path
func InjectMapData(target, subset map[string]interface{}, path string) (map[string]interface{}, error) {
	pathSlice := splitPath(path)
	// Convert the target and subset maps to yaml nodes
	targetNode, err := yaml.FromMap(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create target node from map: %v", err)
	}

	subsetNode, err := yaml.FromMap(subset)
	if err != nil {
		return nil, fmt.Errorf("failed to create subset node from map: %v", err)
	}

	// Get the subset node from target
	targetSubsetNode, err := targetNode.Pipe(yaml.LookupCreate(yaml.MappingNode, pathSlice...))
	if err != nil {
		return nil, fmt.Errorf("error identifying subset node: %v", err)
	}

	// Alternate merge based on custom merge function
	err = mergeYAMLNodes(targetSubsetNode, subsetNode)
	if err != nil {
		return nil, fmt.Errorf("error merging subset into target: %v", err)
	}

	if pathSlice[0] == "" {
		targetNode = targetSubsetNode
	} else {
		if err := targetNode.PipeE(yaml.Lookup(pathSlice[:len(pathSlice)-1]...), yaml.SetField(pathSlice[len(pathSlice)-1], targetSubsetNode)); err != nil {
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

// splitPath splits a path by '.' into a path array
// TODO: This could be a more complicated path: is there a lib function to do this and possibly handle things like [] or escaped '.'
func splitPath(path string) []string {
	// strip leading '.' if present
	if len(path) > 0 && path[0] == '.' {
		path = path[1:]
	}
	return strings.Split(path, ".")
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

package transform

import (
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/utils"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"sigs.k8s.io/kustomize/kyaml/yaml/merge2"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

type ChangeType string

const (
	ChangeTypeAdd    ChangeType = "add"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

type TransformTarget struct {
	RootNode *yaml.RNode
}

func CreateTransformTarget(parent map[string]interface{}) (*TransformTarget, error) {
	// Convert the target to yaml nodes
	node, err := yaml.FromMap(parent)
	if err != nil {
		return nil, fmt.Errorf("failed to create target node from map: %v", err)
	}

	return &TransformTarget{
		RootNode: node,
	}, nil
}

func (t *TransformTarget) ExecuteTransform(path string, cType ChangeType, value string, valueMap map[string]interface{}) (map[string]interface{}, error) {
	rootNodeCopy := t.RootNode.Copy()
	pathSlice, lastSegment, err := CalcPath(path, cType)
	message.Debugf("Path Slice: %v\nLast Item: %s", pathSlice, lastSegment)
	if err != nil {
		return nil, err
	}
	filters, err := BuildFilters(rootNodeCopy, pathSlice)
	if err != nil {
		return nil, fmt.Errorf("error building filters: %v", err)
	}

	switch cType {
	case ChangeTypeAdd, ChangeTypeUpdate:
		var newNode *yaml.RNode

		node, err := rootNodeCopy.Pipe(filters...)
		if err != nil {
			return nil, fmt.Errorf("error finding node in root: %v", err)
		}

		if valueMap != nil {
			newNode, err = getNewNodeFromMap(lastSegment, valueMap)
			if err != nil {
				return nil, fmt.Errorf("error creating new node from map: %v", err)
			}
		} else {
			if len(pathSlice) == 0 {
				return nil, fmt.Errorf("invalid path<>value, cannot set new node as string at root")
			}
			newNode, err = getNewNodeFromString(lastSegment, value)
			if err != nil {
				return nil, fmt.Errorf("error creating new node from string: %v", err)
			}
		}

		if cType == ChangeTypeAdd {
			err := Add(node, newNode)
			if err != nil {
				return nil, err
			}
		} else {
			node, err = Update(node, newNode)
			if err != nil {
				return nil, err
			}
		}

		// Set the node back into the target
		if len(pathSlice) == 0 {
			rootNodeCopy = node
		} else {
			if err := SetNodeAtPath(rootNodeCopy, node, filters, pathSlice); err != nil {
				return nil, fmt.Errorf("error setting merged node back into target: %v", err)
			}
		}

	case ChangeTypeDelete:
		err := Delete(rootNodeCopy, lastSegment, filters)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("invalid transform type: %s", cType)
	}

	// Write node into map[string]interface{}
	var nodeMap map[string]interface{}
	err = rootNodeCopy.YNode().Decode(&nodeMap)
	if err != nil {
		return nil, fmt.Errorf("error decoding root node: %v", err)
	}

	// Update the original
	t.RootNode = rootNodeCopy

	return nodeMap, nil
}

// Add adds the subset to the target at the path, appends to lists
func Add(node, newNode *yaml.RNode) (err error) {
	return mergeYAMLNodes(node, newNode)
}

// Updates existing data at the path, overwrites lists
func Update(node, newNode *yaml.RNode) (*yaml.RNode, error) {
	return merge2.Merge(newNode, node, yaml.MergeOptions{})
}

// Deletes data at the path
func Delete(node *yaml.RNode, lastSegment string, filters []yaml.Filter) (err error) {
	if lastSegment != "" {
		filters = append(filters, yaml.FieldClearer{Name: lastSegment})
		_, err = node.Pipe(filters...)
		if err != nil {
			return fmt.Errorf("error deleting node key: %v", err)
		}
	} else {
		// TODO: If the last segment is a list, we need to delete the last item in the list
		// doesn't appear there's a kyaml filter to help with this...
		return fmt.Errorf("cannot delete a list entry")
	}

	return nil
}

// SetNodeAtPath injects the updated node into rootNode according to the specified path
func SetNodeAtPath(rootNode *yaml.RNode, node *yaml.RNode, filters []yaml.Filter, pathSlice []string) error {
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
			Element: node.Document(),
			Keys:    keys,
			Values:  values,
		})
	} else {
		filters = append(filters[:len(filters)-1], yaml.SetField(lastSegment, node))
	}

	return rootNode.PipeE(filters...)
}

func CalcPath(path string, cType ChangeType) ([]string, string, error) {
	var lastSegment string
	pathSlice := cleanPath(utils.SmarterPathSplitter(path, "."))

	// Path has rules for different change types
	switch cType {
	case ChangeTypeAdd, ChangeTypeUpdate:
		// Add path is the last segment
		if len(pathSlice) > 0 {
			if isFilter(pathSlice[len(pathSlice)-1]) {
				// If the last segment is a filter, the full path will be used as pathSlice
				return pathSlice, "", nil
			} else {
				return pathSlice[:len(pathSlice)-1], pathSlice[len(pathSlice)-1], nil
			}
		}
	case ChangeTypeDelete:
		if len(pathSlice) == 0 {
			return nil, "", fmt.Errorf("invalid path, cannot delete a root node")
		} else {
			if isFilter(pathSlice[len(pathSlice)-1]) {
				// List entry cannot be deleted
				return nil, "", fmt.Errorf("cannot delete a list entry")
			} else {
				return pathSlice[:len(pathSlice)-1], pathSlice[len(pathSlice)-1], nil
			}
		}
	}

	return pathSlice, lastSegment, nil
}

// cleanPath cleans the path slice
func cleanPath(pathSlice []string) []string {
	for i, p := range pathSlice {
		// Remove escaped double quotes
		p = strings.ReplaceAll(p, "\"", "")
		pathSlice[i] = p

		if isFilter(p) {
			// If there's no equal, assume item is a key, NOT a filter
			if !strings.Contains(p, "=") {
				pathSlice[i] = strings.TrimPrefix(strings.TrimSuffix(p, "]"), "[")
			}
		}
	}
	return pathSlice
}

// getNodeFromMap returns the new node to merge from map input
func getNewNodeFromMap(lastSegment string, valueMap map[string]interface{}) (*yaml.RNode, error) {
	if lastSegment != "" {
		valueMap = map[string]interface{}{
			lastSegment: valueMap,
		}
	}

	return yaml.FromMap(valueMap)
}

// getNodeFromString returns the new node to merge from string input
func getNewNodeFromString(lastSegment string, valueStr string) (*yaml.RNode, error) {
	// if the last segment of the pathSlice is a key, set it as the root map key
	if lastSegment != "" {
		return yaml.FromMap(map[string]interface{}{
			lastSegment: valueStr,
		})
	} else {
		return yaml.NewScalarRNode(valueStr), nil
	}
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

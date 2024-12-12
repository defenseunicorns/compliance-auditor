package transform

import (
	"fmt"
	"strconv"

	"sigs.k8s.io/kustomize/kyaml/yaml"
	"sigs.k8s.io/kustomize/kyaml/yaml/merge2"
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

// UpdateRootNode updates the root node of the transform target and returns the updated node
// as a map[string]interface{}
func (t *TransformTarget) UpdateRootNode(node *yaml.RNode) (map[string]interface{}, error) {
	// Write node into map[string]interface{}
	var nodeMap map[string]interface{}
	err := node.YNode().Decode(&nodeMap)
	if err != nil {
		return nil, fmt.Errorf("error decoding root node: %v", err)
	}

	// Update the original
	t.RootNode = node

	return nodeMap, nil
}

func (t *TransformTarget) ExecuteTransform(path string, cType ChangeType, value string, valueMap map[string]interface{}) (map[string]interface{}, error) {
	rootNodeCopy := t.RootNode.Copy()

	pathParts, filters, err := ResolvePathWithFilters(rootNodeCopy, path)
	if err != nil {
		return nil, fmt.Errorf("error resolving path: %v", err)
	}
	operandIdx := len(pathParts) - 1
	if pathParts[operandIdx].Type != PartTypeScalar {
		// If the last segment is not a scalar, the transform will modify a sequence
		operandIdx = -1
	}

	switch cType {
	case ChangeTypeAdd, ChangeTypeUpdate:
		if value != "" {
			// If change type is add or update and value is set -> direct path to set field
			if operandIdx == -1 {
				return nil, fmt.Errorf("invalid combination, cannot set a string value to a sequence")
			}
			if len(pathParts) == 1 && len(filters) == 0 {
				return nil, fmt.Errorf("invalid path, cannot set a string value to root")
			}

			finalNode := yaml.NewScalarRNode(value)
			filters = append(filters[:len(filters)-1], yaml.SetField(pathParts[operandIdx].Value, finalNode))

			err = rootNodeCopy.PipeE(filters...)
			if err != nil {
				return nil, fmt.Errorf("error setting value: %v", err)
			}
		} else {
			// Otherwise, find the node that will be merged into
			node, err := rootNodeCopy.Pipe(filters...)
			if err != nil {
				return nil, fmt.Errorf("error finding node in root: %v", err)
			}

			newNode, err := yaml.FromMap(valueMap)
			if err != nil {
				return nil, fmt.Errorf("error creating new node from map: %v", err)
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

			// Set the merged node back into the target
			if len(filters) == 0 {
				// If root node is being replaced there will be no filters
				rootNodeCopy = node
			} else {
				if err := SetNodeAtPath(rootNodeCopy, node, filters, pathParts, operandIdx); err != nil {
					return nil, fmt.Errorf("error setting merged node back into target: %v", err)
				}
			}
		}

	case ChangeTypeDelete:
		err := Delete(rootNodeCopy, filters, pathParts, operandIdx)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("invalid transform type: %s", cType)
	}

	return t.UpdateRootNode(rootNodeCopy)
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
func Delete(node *yaml.RNode, filters []yaml.Filter, pathParts []PathPart, operandIdx int) (err error) {
	if len(pathParts) == 1 && len(filters) == 0 {
		return fmt.Errorf("invalid path, cannot delete root")
	}

	if operandIdx == -1 {
		// If the last segment is a filter, find the corresponding item in the list and remove
		// Length of pathParts should be at least 2 to find the sequence node key
		if len(pathParts) < 2 || len(filters) < 2 {
			return fmt.Errorf("invalid path length for index operation")
		}
		sequenceKey := pathParts[len(pathParts)-2].Value
		lastPart := pathParts[len(pathParts)-1]

		// Get all elements, omit the node that matches, then set the sequence back in...
		parentSeq, err := node.Pipe(filters[:len(filters)-1]...)
		if err != nil {
			return err
		}
		if parentSeq == nil {
			return fmt.Errorf("parent sequence is nil")
		}

		// Get the matching node
		matchedNode, err := node.Pipe(filters...)
		if err != nil {
			return err
		}

		if parentSeq.YNode().Kind == yaml.SequenceNode {
			// Get all nodes in the sequence
			allSeqNodes, err := parentSeq.Elements()
			if err != nil {
				return err
			}

			newSeqNode := &yaml.Node{
				Kind:    yaml.SequenceNode,
				Content: make([]*yaml.Node, 0),
			}
			for i, seqNode := range allSeqNodes {
				if lastPart.Type == PartTypeIndex {
					// Convert the index to a usable int
					var idx int
					if lastPart.Value == "-" {
						idx = len(allSeqNodes) - 1
					} else {
						idx, err = strconv.Atoi(lastPart.Value)
						if err != nil {
							return err
						}
					}

					// If the index is not the same, append the node to the new sequence
					if i != idx {
						newSeqNode.Content = append(newSeqNode.Content, seqNode.YNode())
					}
				} else {
					// This is a bit of hack, might be a better way to do this
					matchedBytes, err := matchedNode.MarshalJSON()
					if err != nil {
						return err
					}
					seqBytes, err := seqNode.MarshalJSON()
					if err != nil {
						return err
					}

					if string(matchedBytes) != string(seqBytes) {
						newSeqNode.Content = append(newSeqNode.Content, seqNode.YNode())
					}
				}
			}

			filters = append(filters[:len(filters)-2], yaml.FieldSetter{
				Name:  sequenceKey,
				Value: yaml.NewRNode(newSeqNode),
			})
			_, err = node.Pipe(filters...)
			if err != nil {
				return fmt.Errorf("error deleting node key: %v", err)
			}

		} else {
			return fmt.Errorf("expected sequence node, but got %v", parentSeq.YNode().Kind)
		}

	} else {
		filters = append(filters[:len(filters)-1], yaml.FieldClearer{Name: pathParts[operandIdx].Value})
		_, err = node.Pipe(filters...)
		if err != nil {
			return fmt.Errorf("error deleting node key: %v", err)
		}
	}

	return nil
}

// SetNodeAtPath injects the updated node into rootNode according to the specified path and final node type
func SetNodeAtPath(rootNode *yaml.RNode, node *yaml.RNode, filters []yaml.Filter, pathParts []PathPart, finalItemIdx int) error {
	if len(filters) == 0 {
		return fmt.Errorf("must have at least one filter")
	}

	if finalItemIdx == -1 {
		lastPart := pathParts[len(pathParts)-1]
		if lastPart.Type == PartTypeSelector {
			// Last item is a filter
			selectorParts, err := extractSelector(lastPart.Value)
			if err != nil {
				return err
			}

			keys := make([]string, 0)
			values := make([]string, 0)
			for _, part := range selectorParts {
				if isComposite(lastPart.Value) {
					// This shouldn't happen
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

		} else { // Last item is an index
			// Length of pathParts should be >= 2 to find the sequence node key
			if len(pathParts) < 2 || len(filters) < 2 {
				return fmt.Errorf("invalid path length for index operation")
			}
			sequenceKey := pathParts[len(pathParts)-2].Value

			// Get all elements, update the node that matches, then set the sequence back in...
			parentSeq, err := rootNode.Pipe(filters[:len(filters)-1]...)
			if err != nil {
				return err
			}
			if parentSeq == nil {
				return fmt.Errorf("parent sequence is nil")
			}

			if parentSeq.YNode().Kind == yaml.SequenceNode {
				// Get all nodes in the sequence
				allSeqNodes, err := parentSeq.Elements()
				if err != nil {
					return err
				}

				// Convert the index to a usable int
				var idx int
				if lastPart.Value == "-" {
					idx = len(allSeqNodes) - 1
				} else {
					idx, err = strconv.Atoi(lastPart.Value)
					if err != nil {
						return err
					}
				}

				newSeqNode := &yaml.Node{
					Kind:    yaml.SequenceNode,
					Content: make([]*yaml.Node, 0),
				}
				for i, seqNode := range allSeqNodes {
					if i == idx {
						newSeqNode.Content = append(newSeqNode.Content, node.YNode())
					} else {
						newSeqNode.Content = append(newSeqNode.Content, seqNode.YNode())
					}
				}

				filters = append(filters[:len(filters)-2], yaml.FieldSetter{
					Name:  sequenceKey,
					Value: yaml.NewRNode(newSeqNode),
				})
			}
		}
	} else {
		// Replace the last filter (which should be identifying the node) with a SetField filter
		filters = append(filters[:len(filters)-1], yaml.SetField(pathParts[finalItemIdx].Value, node))
	}

	return rootNode.PipeE(filters...)
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

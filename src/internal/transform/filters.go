package transform

import (
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type selectorPart struct {
	key   string
	value string
}

// ResolvePathWithFilters converts a path to its individual parts which define the type of part
// along with the appropriate kyaml filters to apply to the path
func ResolvePathWithFilters(targetNode *yaml.RNode, path string) ([]PathPart, []yaml.Filter, error) {
	pathParts := PathToParts(path) // Will always return at least one item

	filters := make([]yaml.Filter, 0)
	for i, part := range pathParts {
		if part.Type == PartTypeSelector {
			selectorParts, err := extractSelector(part.Value)
			if err != nil {
				return nil, nil, err
			}

			if len(selectorParts) > 1 || isComposite(selectorParts[0].key) {
				index, err := returnIndexFromComplexSelectors(targetNode, filters, selectorParts)
				if err != nil {
					return nil, nil, err
				}

				if index == -1 {
					return nil, nil, fmt.Errorf("composite path not found")
				} else {
					filters = append(filters, yaml.GetElementByIndex(index))
					// Update the partPath with the index instead of composite key
					pathParts[i].Value = fmt.Sprintf("%d", index)
					pathParts[i].Type = PartTypeIndex
				}
			} else {
				filters = append(filters, yaml.MatchElement(trimDoubleQuotes(selectorParts[0].key), selectorParts[0].value))
			}
		} else {
			if part.Value != "" {
				filters = append(filters, yaml.Lookup(trimDoubleQuotes(part.Value)))
			}
		}

	}

	// Validate that the filters are valid for the target node
	err := targetNode.PipeE(filters...)
	if err != nil {
		return nil, nil, err
	}

	return pathParts, filters, nil
}

// isComposite checks if a string is a composite selector, e.g., metadata.name
// input is not composite if encapsulated in quotes, e.g., "metadata.name" -> it's the value of the key
func isComposite(input string) bool {
	if input[0] == '"' && input[len(input)-1] == '"' {
		return false
	}

	keys := strings.Split(input, ".")
	return len(keys) > 1
}

// buildCompositeFilters creates a yaml.Filter slice for a composite selector
// e.g., [metadata.namespace=foo]
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

// extractSelector extracts the selector parts from a string
// e.g., key1=value1,key2=value2; composite.key=value; val.key.test=bar
func extractSelector(item string) ([]selectorPart, error) {
	items := strings.Split(item, ",")
	if len(items) == 0 {
		return []selectorPart{}, fmt.Errorf("filter is empty")
	}

	selelctorParts := make([]selectorPart, 0, len(items))
	for _, i := range items {
		if !strings.Contains(i, "=") {
			return []selectorPart{}, fmt.Errorf("filter is not in the correct format")
		}
		selelctorParts = append(selelctorParts, selectorPart{
			key:   strings.SplitN(i, "=", 2)[0],
			value: strings.SplitN(i, "=", 2)[1],
		})
	}

	return selelctorParts, nil
}

// returnIndexFromComplexSelectors returns the index of the node that matches the filterParts
// e.g., [key1=value1,key2=value2], [composite.key=value], [val.key.test=bar]
func returnIndexFromComplexSelectors(targetNode *yaml.RNode, parentFilters []yaml.Filter, selectorParts []selectorPart) (int, error) {
	if targetNode == nil {
		return -1, fmt.Errorf("root node cannot be nil")
	}

	index := -1

	parentNode, err := targetNode.Pipe(parentFilters...)
	if err != nil {
		return index, err
	}

	if parentNode == nil {
		return index, fmt.Errorf("parent node is not found for filters: %v", parentFilters)
	}

	if parentNode.YNode().Kind == yaml.SequenceNode {
		nodes, err := parentNode.Elements()
		if err != nil {
			return index, err
		}
		for i, node := range nodes {
			if nodeMatchesAllFilters(node, selectorParts) {
				index = i
				break
			}
		}
	} else {
		return index, fmt.Errorf("expected sequence node, but got %v", parentNode.YNode().Kind)
	}

	return index, nil
}

func nodeMatchesAllFilters(node *yaml.RNode, selectorParts []selectorPart) bool {
	for _, part := range selectorParts {
		if isComposite(part.key) {
			compositeFilters := buildCompositeFilters(part.key, part.value)
			n, err := node.Pipe(compositeFilters...)
			if err != nil || n == nil {
				return false
			}
		} else {
			n, err := node.Pipe(yaml.MatchField(trimDoubleQuotes(part.key), part.value))
			if err != nil || n == nil {
				return false
			}
		}
	}
	return true
}

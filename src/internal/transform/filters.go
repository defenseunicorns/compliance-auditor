package transform

import (
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type filterParts struct {
	key   string
	value string
}

// BuildFilters builds kyaml filters from a pathSlice
func BuildFilters(targetNode *yaml.RNode, pathSlice []string) ([]yaml.Filter, error) {
	if targetNode == nil {
		return nil, fmt.Errorf("root node is nil")
	}

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

// Helper to determine if item in path is a filter
func isFilter(item string) bool {
	// check if first and last char are [ and ]
	return strings.HasPrefix(item, "[") && strings.HasSuffix(item, "]")
}

// isComposite checks if a string is a composite string, e.g., metadata.name
func isComposite(input string) bool {
	keys := strings.Split(input, ".")
	return len(keys) > 1
}

// buildCompositeFilters creates a yaml.Filter slice for a composite key
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

// returnIndexFromComplexFilters returns the index of the node that matches the filterParts
// e.g., [key1=value1,key2=value2], [composite.key=value], [val.key.test=bar]
func returnIndexFromComplexFilters(targetNode *yaml.RNode, parentFilters []yaml.Filter, filterParts []filterParts) (int, error) {
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

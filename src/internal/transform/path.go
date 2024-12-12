package transform

import (
	"regexp"
	"strconv"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/utils"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

type PartType int

const (
	PartTypeMap      PartType = iota // e.g., a in a.b
	PartTypeSequence                 // e.g., a in a[b=c]
	PartTypeScalar                   // e.g., b in a.b
	PartTypeSelector                 // e.g., [a=b]
	PartTypeIndex                    // e.g., [0]
)

type PathPart struct {
	Type  PartType
	Value string
}

// PathToParts converts the path string into a slice of pathParts
func PathToParts(path string) []PathPart {
	pathSlice := utils.SmarterPathSplitter(normalizePath(path), ".")
	message.Debug("Path Slice: %v\n", pathSlice) // Helpful for understanding issues with how a path is parsed

	return makePathParts(pathSlice)
}

// Normalize the path to kyaml syntax by inserting any missing "." before "["
// e.g., "foo[bar=baz]" -> "foo.[bar=baz]"
func normalizePath(path string) string {
	var builder strings.Builder

	// Iterate over the string
	for i, char := range path {
		// If the current character is `[`
		if char == '[' {
			// Check if it's at the start of the string or if the preceding character is not `.`
			if i == 0 || path[i-1] != '.' {
				builder.WriteRune('.') // Inject a `.`
			}
		}
		// Append the current character
		builder.WriteRune(char)
	}

	// Get the modified string
	result := builder.String()

	return result
}

// makePathParts creates the pathParts from the pathSlice
func makePathParts(pathSlice []string) []PathPart {
	pathParts := make([]PathPart, 0, len(pathSlice))

	for i, p := range pathSlice {
		p = cleanPart(p)
		currentPartType := getPartType(p)

		var pathPart PathPart

		// If the current part is a scalar, look ahead to see if it's a map or sequence
		if currentPartType == PartTypeScalar && i < len(pathSlice)-1 {
			nextPart := cleanPart(pathSlice[i+1])
			pathPart = pathPartFromLookAhead(trimDoubleQuotes(p), getPartType(nextPart))
		} else {
			// Calculate the pathPart from the current element
			pathPart = PathPart{Type: currentPartType, Value: trimDoubleQuotes(p)}
		}

		pathParts = append(pathParts, pathPart)
	}

	return pathParts
}

// Removes leading and trailing brackets
// TODO: revisit if we need to clean escaped quotes in other spots, esp if supporting json
func cleanPart(part string) string {
	// Trim any leading or trailing brackets
	part = strings.TrimPrefix(strings.TrimSuffix(part, "]"), "[")

	// Trim any leading or trailing escaped quotes
	if strings.HasPrefix(part, `\"`) && strings.HasSuffix(part, `\"`) {
		part = strings.TrimPrefix(strings.TrimSuffix(part, `\"`), `\"`)
		part = `"` + part + `"`
	}

	return part
}

// trimDoubleQuotes removes leading and trailing double quotes, if present
func trimDoubleQuotes(part string) string {
	if strings.HasPrefix(part, `"`) && strings.HasSuffix(part, `"`) {
		part = strings.TrimPrefix(strings.TrimSuffix(part, `"`), `"`)
	}
	return part
}

func pathPartFromLookAhead(current string, nextType PartType) PathPart {
	switch nextType {
	case PartTypeSelector, PartTypeIndex:
		return PathPart{Type: PartTypeSequence, Value: current}
	case PartTypeScalar:
		fallthrough
	default:
		return PathPart{Type: PartTypeMap, Value: current}
	}
}

func getPartType(part string) PartType {
	if isFilter(part) {
		// Check if the part is a filter (e.g. [a=b] or [a.b=c] or ["a.b"=c])
		return PartTypeSelector
	} else if isListIndex(part) {
		// Check if the part is a list index (e.g. 0 or -)
		return PartTypeIndex
	} else {
		// Anything else is a scalar
		return PartTypeScalar
	}
}

// Helper to determine if part is a filter
func isFilter(item string) bool {
	// check if part contains "=" that is not escaped or encapsulted in quotes
	// Note - "some.key"=value will return filter, but "some.key"="value" will NOT (even though it is), so the latter should not be used
	re := regexp.MustCompile(`".*?=.*?"`)
	return strings.Contains(item, "=") && !re.MatchString(item)
}

// Helper to determine if part is a list index
func isListIndex(item string) bool {
	if item == "-" {
		return true
	}

	_, err := strconv.Atoi(item)
	return err == nil
}

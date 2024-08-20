package oscal

import (
	"regexp"
	"strconv"
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

const (
	LULA_NAMESPACE = "https://docs.lula.dev/oscal/ns"
	LULA_KEYWORD   = "lula"
)

// Update legacy_namespaces when namespace URL (LULA_NAMESPACE) changes to ensure backwards compatibility
var legacy_namespaces = []string{"https://docs.lula.dev/ns"}

// UpdateProps updates a property in a slice of properties or adds if not exists
func UpdateProps(name string, namespace string, value string, props *[]oscalTypes_1_1_2.Property) {

	for index, prop := range *props {
		found, propNamespace := checkOrUpdateNamespace(prop.Ns, namespace)
		if prop.Name == name && found {
			prop.Value = value
			prop.Ns = propNamespace
			(*props)[index] = prop
			return
		}
	}
	// Prop does not exist
	prop := oscalTypes_1_1_2.Property{
		Ns:    namespace,
		Name:  name,
		Value: value,
	}

	*props = append(*props, prop)
}

func GetProp(name string, namespace string, props *[]oscalTypes_1_1_2.Property) (bool, string) {

	if props == nil {
		return false, ""
	}

	for _, prop := range *props {
		found, _ := checkOrUpdateNamespace(prop.Ns, namespace)
		if prop.Name == name && found {
			return true, prop.Value
		}
	}
	return false, ""
}

func checkOrUpdateNamespace(propNamespace, namespace string) (bool, string) {
	// if namespace doesn't contain lula, check namespace == propNamespace
	if !strings.Contains(propNamespace, LULA_KEYWORD) {
		return namespace == propNamespace, propNamespace
	}

	for _, ns := range legacy_namespaces {
		if namespace == ns {
			return true, LULA_NAMESPACE
		}
	}
	return namespace == LULA_NAMESPACE, LULA_NAMESPACE
}

// CompareControls compares two control titles, handling both XX-##.## formats and regular strings.
// true sorts a before b; false sorts b before a
func CompareControls(a, b string) bool {
	// Define a regex to match the XX-##.## format
	nistFormat := regexp.MustCompile(`(?i)^[a-z]{2}-\d+(\.\d+)?$`)

	// Check if both strings match the XX-##.## format
	isANistFormat := nistFormat.MatchString(a)
	isBNistFormat := nistFormat.MatchString(b)

	// If both are in XX-##.## format, apply the custom comparison logic
	if isANistFormat && isBNistFormat {
		return compareNistFormat(a, b)
	}

	// If neither are in XX-##.## format, use simple lexicographical comparison
	if !isANistFormat && !isBNistFormat {
		return a < b
	}

	// If only one is in XX-##.## format, treat it as "less than" the regular string
	return !isANistFormat
}

// compareNistFormat handles the comparison for strings in the XX-##.## format.
func compareNistFormat(a, b string) bool {
	// Split the strings by "-"
	splitA := strings.Split(a, "-")
	splitB := strings.Split(b, "-")

	// Compare the alphabetic part first
	if splitA[0] != splitB[0] {
		return splitA[0] < splitB[0]
	}

	// Compare the numeric part before the dot (.)
	numA, _ := strconv.Atoi(strings.Split(splitA[1], ".")[0])
	numB, _ := strconv.Atoi(strings.Split(splitB[1], ".")[0])

	if numA != numB {
		return numA < numB
	}

	// Compare the numeric part after the dot (.) if exists
	if len(strings.Split(splitA[1], ".")) > 1 && len(strings.Split(splitB[1], ".")) > 1 {
		subNumA, _ := strconv.Atoi(strings.Split(splitA[1], ".")[1])
		subNumB, _ := strconv.Atoi(strings.Split(splitB[1], ".")[1])
		return subNumA < subNumB
	}

	// Handle cases where only one has a sub-number
	if len(strings.Split(splitA[1], ".")) > 1 {
		return false
	}
	if len(strings.Split(splitB[1], ".")) > 1 {
		return true
	}

	return false
}

package oscal

import (
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

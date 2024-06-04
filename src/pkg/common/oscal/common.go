package oscal

import (
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
)

// UpdateProps updates a property in a slice of properties or adds if not exists
func UpdateProps(name string, namespace string, value string, props *[]oscalTypes_1_1_2.Property) {

	for index, prop := range *props {
		if prop.Name == name && prop.Ns == namespace {
			prop.Value = value
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

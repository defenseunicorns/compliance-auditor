package oscal

import (
	"fmt"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

// NewCatalog creates a new catalog object from the given data.
func NewCatalog(data []byte) (catalog *oscalTypes.Catalog, err error) {
	var oscalModels oscalTypes.OscalModels

	// validate the catalog
	err = multiModelValidate(data)
	if err != nil {
		return catalog, err
	}

	// unmarshal the catalog
	// yaml.v3 unmarshal handles both json and yaml
	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		message.Debugf("Error marshalling yaml: %s\n", err.Error())
		return catalog, err

	}

	return oscalModels.Catalog, nil
}

// ResolveCatalogControls resolves all controls in the provided catalog
func ResolveCatalogControls(catalog *oscalTypes.Catalog, include, exclude []string) (map[string]oscalTypes.Control, error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog is nil")
	}

	if len(include) != 0 && len(exclude) != 0 {
		return nil, fmt.Errorf("include and exclude cannot be used together")
	}

	controlMap := make(map[string]oscalTypes.Control)
	var err error

	// Begin recursive control search among groups and controls
	if catalog.Groups != nil {
		controlMap, err = searchGroups(catalog.Groups, controlMap, include, exclude)
		if err != nil {
			return nil, err
		}
	}

	if catalog.Controls != nil {
		controlMap, err = searchControls(catalog.Controls, controlMap, include, exclude)
		if err != nil {
			return nil, err
		}
	}

	return controlMap, nil
}

// getControlRemarks gets the control-remarks from the provided control, as determined by the targetRemarks
func getControlRemarks(control *oscalTypes.Control, targetRemarks []string) (string, error) {
	var controlRemarks string
	paramMap := make(map[string]parameter)

	if control == nil {
		return "", fmt.Errorf("control is nil")
	}

	if control.Params != nil {
		for _, param := range *control.Params {

			if param.Select == nil {
				paramMap[param.ID] = parameter{
					ID:    param.ID,
					Label: param.Label,
				}
			} else {
				sel := *param.Select
				paramMap[param.ID] = parameter{
					ID: param.ID,
					Select: &selection{
						HowMany: sel.HowMany,
						Choice:  *sel.Choice,
					},
				}
			}
		}
	} else {
		message.Debugf("No parameters (control.Params) found for %s", control.ID)
	}

	if control.Parts != nil {
		for _, part := range *control.Parts {
			if contains(targetRemarks, part.Name) {
				controlRemarks += fmt.Sprintf("%s:\n", strings.ToTitle(part.Name))
				if part.Prose != "" && strings.Contains(part.Prose, "{{ insert: param,") {
					controlRemarks += replaceParams(part.Prose, paramMap, false)
				} else {
					controlRemarks += part.Prose
				}
				if part.Parts != nil {
					controlRemarks += addPart(part.Parts, paramMap, 0)
				}
			}
		}
	}

	// assemble implemented-requirements object
	return controlRemarks, nil
}

func searchGroups(groups *[]oscalTypes.Group, controlMap map[string]oscalTypes.Control, include, exclude []string) (map[string]oscalTypes.Control, error) {
	var err error

	for _, group := range *groups {
		if group.Groups != nil {
			controlMap, err = searchGroups(group.Groups, controlMap, include, exclude)
			if err != nil {
				return nil, err
			}
		}
		if group.Controls != nil {
			controlMap, err = searchControls(group.Controls, controlMap, include, exclude)
			if err != nil {
				return nil, err
			}
		}
	}
	return controlMap, nil
}

func searchControls(controls *[]oscalTypes.Control, controlMap map[string]oscalTypes.Control, include, exclude []string) (map[string]oscalTypes.Control, error) {
	var err error

	for _, control := range *controls {
		// Add control if specified
		if AddControl(control.ID, include, exclude) {
			// If the control is not already in the map, add it
			if _, ok := controlMap[control.ID]; !ok {
				controlMap[control.ID] = control
			}
		}

		// Check all child controls
		if control.Controls != nil {
			controlMap, err = searchControls(control.Controls, controlMap, include, exclude)
			if err != nil {
				return nil, err
			}
		}
	}
	return controlMap, nil
}

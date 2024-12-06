package reporting

import (
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

// Split the default controlMap into framework and source maps for further processing
func SplitControlMap(controlMap map[string][]oscalTypes.ControlImplementationSet) (sourceMap map[string]map[string]int, frameworkMap map[string]map[string]int) {
	sourceMap = make(map[string]map[string]int)
	frameworkMap = make(map[string]map[string]int)

	for key, implementations := range controlMap {
		for _, controlImplementation := range implementations {
			status, framework := oscal.GetProp("framework", oscal.LULA_NAMESPACE, controlImplementation.Props)
			if status {
				// if these are the same - we need to de-duplicate
				if key == framework {
					if _, exists := frameworkMap[framework]; !exists {
						frameworkMap[framework] = make(map[string]int)
					}
					for _, implementedReq := range controlImplementation.ImplementedRequirements {
						controlID := implementedReq.ControlId
						frameworkMap[framework][controlID]++
					}
				} else {
					if _, exists := sourceMap[key]; !exists {
						sourceMap[key] = make(map[string]int)
					}
					for _, implementedReq := range controlImplementation.ImplementedRequirements {
						controlID := implementedReq.ControlId
						sourceMap[key][controlID]++
					}
				}
			} else {
				if _, exists := sourceMap[key]; !exists {
					sourceMap[key] = make(map[string]int)
				}
				for _, implementedReq := range controlImplementation.ImplementedRequirements {
					controlID := implementedReq.ControlId
					sourceMap[key][controlID]++
				}
			}
		}
	}

	return sourceMap, frameworkMap
}

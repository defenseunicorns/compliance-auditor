package reporting

import (
	"os"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/message"
)

// Helper to determine if the controlMap source is a URL
func isURL(str string) bool {
	_, err := network.ParseUrl(str)
	return err == nil
}

func fetchOrReadFile(source string) ([]byte, error) {
	if isURL(source) {
		spinner := message.NewProgressSpinner("Fetching data from URL: %s", source)
		defer spinner.Stop()
		data, err := network.Fetch(source)
		if err != nil {
			spinner.Fatalf(err, "failed to fetch data from URL")
		}
		spinner.Success()
		return data, nil
	}
	spinner := message.NewProgressSpinner("Reading file: %s", source)
	defer spinner.Stop()
	data, err := os.ReadFile(source)
	if err != nil {
		spinner.Fatalf(err, "failed to read file")
	}
	spinner.Success()
	return data, nil
}

// Split the default controlMap into framework and source maps for further processing
func SplitControlMap(controlMap map[string][]oscalTypes_1_1_2.ControlImplementationSet) (sourceMap map[string]map[string]int, frameworkMap map[string]map[string]int) {
	sourceMap = make(map[string]map[string]int)
	frameworkMap = make(map[string]map[string]int)

	for key, implementations := range controlMap {
		if isURL(key) {
			if _, exists := sourceMap[key]; !exists {
				sourceMap[key] = make(map[string]int)
			}
			for _, controlImplementation := range implementations {
				for _, implementedReq := range controlImplementation.ImplementedRequirements {
					controlID := implementedReq.ControlId
					sourceMap[key][controlID]++
				}
			}
		} else {
			if _, exists := frameworkMap[key]; !exists {
				frameworkMap[key] = make(map[string]int)
			}
			for _, controlImplementation := range implementations {
				for _, implementedReq := range controlImplementation.ImplementedRequirements {
					controlID := implementedReq.ControlId
					frameworkMap[key][controlID]++
				}
			}
		}
	}

	return sourceMap, frameworkMap
}

package oscal

import (
	"encoding/json"
	"fmt"
	"strings"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"gopkg.in/yaml.v3"
)

func NewCatalog(source string, data []byte) (catalog oscalTypes_1_1_2.Catalog, err error) {
	var oscalModels oscalTypes_1_1_2.OscalModels
	if strings.HasSuffix(source, ".yaml") {
		err = yaml.Unmarshal(data, &oscalModels)
		if err != nil {
			message.Debugf("Error marshalling yaml: %s\n", err.Error())
			return catalog, err
		}
	} else if strings.HasSuffix(source, ".json") {
		err = json.Unmarshal(data, &oscalModels)
		if err != nil {
			message.Debugf("Error marshalling json: %s\n", err.Error())
			return catalog, err
		}
	} else {
		return catalog, fmt.Errorf("unsupported file type: %s", source)
	}

	return *oscalModels.Catalog, nil
}

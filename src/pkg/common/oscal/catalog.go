package oscal

import (
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"gopkg.in/yaml.v3"
)

func NewCatalog(data []byte) (catalog oscalTypes_1_1_2.Catalog, err error) {
	var oscalModels oscalTypes_1_1_2.OscalModels

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		message.Debugf("Error marshalling yaml: %s\n", err.Error())
		return catalog, err
	}

	return *oscalModels.Catalog, nil
}

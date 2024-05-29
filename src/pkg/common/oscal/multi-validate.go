package oscal

import (
	"errors"

	"github.com/defenseunicorns/go-oscal/src/pkg/model"
	"github.com/defenseunicorns/go-oscal/src/pkg/validation"
)

func multiModelValidate(data []byte) (err error) {
	jsonMap, err := model.CoerceToJsonMap(data)
	if err != nil {
		return err
	}

	if len(jsonMap) == 0 {
		return errors.New("no models found")
	}

	for key, value := range jsonMap {
		jsonModel := make(map[string]interface{})
		jsonModel[key] = value
		validator, err := validation.NewValidator(jsonModel)
		if err != nil {
			return err
		}

		err = validator.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

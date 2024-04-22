package validationstore

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
	"github.com/defenseunicorns/lula/src/types"
)

const UUID_PREFIX = "#"
const WILDCARD = "*"
const YAML_DELIMITER = "---"

type ValidationStore struct {
	backMatterMap map[string]string
	validationMap types.LulaValidationMap
	hrefIdMap     map[string][]string
}

// NewValidationStore creates a new validation store
func NewValidationStore() *ValidationStore {
	return &ValidationStore{
		backMatterMap: make(map[string]string),
		validationMap: make(types.LulaValidationMap),
		hrefIdMap:     make(map[string][]string),
	}
}

// NewValidationStoreFromBackMatter creates a new validation store from a back matter
func NewValidationStoreFromBackMatter(backMatter oscalTypes_1_1_2.BackMatter) *ValidationStore {
	return &ValidationStore{
		backMatterMap: oscal.BackMatterToMap(backMatter),
		validationMap: make(types.LulaValidationMap),
		hrefIdMap:     make(map[string][]string),
	}
}

// AddValidation adds a validation to the store
func (v *ValidationStore) AddValidation(validation *common.Validation) (id string, err error) {
	if validation.Metadata.UUID == "" {
		validation.Metadata.UUID = uuid.NewUUID()
	}

	v.validationMap[validation.Metadata.UUID], err = validation.ToLulaValidation()

	if err != nil {
		return "", err
	}

	return validation.Metadata.UUID, nil

}

// GetLulaValidation gets the LulaValidation from the store
func (v *ValidationStore) GetLulaValidation(id string) (validation *types.LulaValidation, err error) {
	trimmedId := TrimIdPrefix(id)

	if validation, ok := v.validationMap[trimmedId]; ok {
		return &validation, nil
	}

	if validationString, ok := v.backMatterMap[trimmedId]; ok {
		lulaValidation, err := common.ValidationFromString(validationString)
		if err != nil {
			return nil, err
		}
		v.validationMap[trimmedId] = lulaValidation
		return &lulaValidation, nil
	}

	return validation, fmt.Errorf("validation #%s not found", trimmedId)
}

// SetHrefIds sets the validation ids for a given href
func (v *ValidationStore) SetHrefIds(href string, ids []string) {
	v.hrefIdMap[href] = ids
}

// GetHrefIds gets the validation ids for a given href
func (v *ValidationStore) GetHrefIds(href string) (ids []string, err error) {
	if ids, ok := v.hrefIdMap[href]; ok {
		return ids, nil
	}
	return nil, fmt.Errorf("href #%s not found", href)
}

// AddFromLink adds a validation from a link
func (v *ValidationStore) AddFromLink(link oscalTypes_1_1_2.Link) (ids []string, err error) {
	id := link.Href

	// If the resource fragment is not a wildcard, trim the prefix from the resource fragment
	if link.ResourceFragment != WILDCARD && link.ResourceFragment != "" {
		id = TrimIdPrefix(link.ResourceFragment)
	}

	// If the id is a uuid and the lula validation exists, return the id
	if _, err := v.GetLulaValidation(id); err == nil {
		return []string{id}, err
	}

	// If the id is a url and has been fetched before, return the ids
	if ids, err := v.GetHrefIds(id); err == nil {
		return ids, nil
	}

	// If the id is a url and has not been fetched before, fetch and add to the store
	ids, err = v.fetchFromRemoteLink(link)
	if err != nil {
		return ids, err
	}

	return ids, nil
}

// fetchFromRemoteLink adds a validation from a remote source
func (v *ValidationStore) fetchFromRemoteLink(link oscalTypes_1_1_2.Link) (ids []string, err error) {
	wantedId := TrimIdPrefix(link.ResourceFragment)

	validationBytes, err := network.Fetch(link.Href)
	if err != nil {
		return ids, err
	}

	validationBytesArr := bytes.Split(validationBytes, []byte(YAML_DELIMITER))
	isSingleValidation := len(validationBytesArr) == 1

	for _, validationBytes := range validationBytesArr {
		var validation common.Validation
		if err = validation.UnmarshalYaml(validationBytes); err != nil {
			return ids, err
		}
		// If the validation does not have a UUID, create a new one
		if validation.Metadata.UUID == "" {
			validation.Metadata.UUID = uuid.NewUUID()
		}

		// Add the validation to the store
		id, err := v.AddValidation(&validation)
		if err != nil {
			return ids, err
		}

		// If the wanted id is the id, the id is a wildcard, or there is only one validation, add the id to the ids
		if wantedId == id || wantedId == WILDCARD || isSingleValidation {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return ids, fmt.Errorf("no validations found for %s", link.Href)
	} else {
		v.SetHrefIds(link.Href, ids)
	}

	return ids, nil
}

// TrimIdPrefix trims the id prefix from the given id
func TrimIdPrefix(id string) string {
	return strings.TrimPrefix(id, UUID_PREFIX)
}

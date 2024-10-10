package oscal

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"gopkg.in/yaml.v3"
)

type Profile struct {
	Model *oscalTypes.Profile
}

func NewProfile() *Profile {
	var profile Profile
	profile.Model = nil
	return &profile
}

func (p *Profile) GetType() string {
	return "profile"
}

func (p *Profile) GetCompleteModel() *oscalTypes.OscalModels {
	return &oscalTypes.OscalModels{
		Profile: p.Model,
	}
}

func (p *Profile) MakeDeterministic() error {
	if p.Model == nil {
		return fmt.Errorf("cannot make nil model deterministic")
	} else {
		// sort the import items by source string
		importItems := p.Model.Imports

		sort.Slice(importItems, func(i, j int) bool {
			return importItems[i].Href < importItems[j].Href
		})

		// Does not handle pattern matching
		for _, item := range importItems {

			// Shouldn't be both but this functionality isn't the gate
			if item.IncludeControls != nil {
				includeControls := *item.IncludeControls
				for _, includeControl := range includeControls {
					includes := *includeControl.WithIds
					slices.SortStableFunc(includes, func(a, b string) int {
						return CompareControlsInt(a, b)
					})
					includeControl.WithIds = &includes
				}
			}

			if item.ExcludeControls != nil {
				excludeControls := *item.ExcludeControls
				for _, excludeControl := range excludeControls {
					exclude := *excludeControl.WithIds
					slices.SortStableFunc(exclude, func(a, b string) int {
						return CompareControlsInt(a, b)
					})
					excludeControl.WithIds = &exclude
				}
			}

		}

		// sort backmatter
		if p.Model.BackMatter != nil {
			backmatter := *p.Model.BackMatter
			sortBackMatter(&backmatter)
			p.Model.BackMatter = &backmatter
		}
	}

	return nil
}

func (p *Profile) HandleExisting(filepath string) error {
	exists, err := common.CheckFileExists(filepath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("output File %s currently exist - cannot merge artifacts", filepath)
	} else {
		return nil
	}
}

// Create a new profile model
func (p *Profile) NewModel(data []byte) error {

	var oscalModels oscalTypes.OscalModels

	err := multiModelValidate(data)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		return err
	}

	p.Model = oscalModels.Profile

	return nil
}

func GenerateProfile(command string, source string, include []string, exclude []string) (*Profile, error) {

	// Create the OSCAL profile type model for use and later assignment to the oscal.Profile implementation
	var model oscalTypes.Profile

	// Single time used for all time related fields
	rfc3339Time := time.Now()

	// Always create a new UUID for the assessment results (for now)
	model.UUID = uuid.NewUUID()

	// Creation of the generation prop
	props := []oscalTypes.Property{
		{
			Name:  "generation",
			Ns:    LULA_NAMESPACE,
			Value: command,
		},
	}

	// Create metadata object with requires fields and a few extras
	// Adding props to metadata as it is less available within the model
	model.Metadata = oscalTypes.Metadata{
		Title:        "Profile",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "Profile generated from Lula",
		Published:    &rfc3339Time,
		LastModified: rfc3339Time,
		Props:        &props,
	}

	// Include would include the specified controls and exclude the rest
	// Exclude would exclude the specified controls and include the rest
	// Both doesn't make sense - TODO: Need to validate what OSCAL supports here
	includedControls := []oscalTypes.SelectControlById{
		{
			WithIds: &include,
		},
	}

	excludedControls := []oscalTypes.SelectControlById{
		{
			WithIds: &exclude,
		},
	}

	importItem := oscalTypes.Import{
		Href: source,
	}

	// Handle the inclusion of both before passed into this function
	if len(include) > 0 {
		importItem.IncludeControls = &includedControls
	}

	if len(exclude) > 0 {
		importItem.ExcludeControls = &excludedControls
	}

	model.Imports = []oscalTypes.Import{
		importItem,
	}

	// Static allocation of the merge setting until other use-cases are identified
	model.Merge = &oscalTypes.Merge{
		AsIs: true,
	}

	var profile Profile

	profile.Model = &model

	return &profile, nil

}

package oscal

import (
	"fmt"
	"sort"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"gopkg.in/yaml.v3"
)

type Profile struct {
	Model     *oscalTypes.Profile
	ModelType string
}

func (p Profile) GetType() string {
	return "profile"
}

func (p Profile) GetCompleteModel() *oscalTypes.OscalModels {
	return &oscalTypes.OscalModels{
		Profile: p.Model,
	}
}

func (p Profile) MakeDeterministic() {

	// sort the import items by source string
	importItems := p.Model.Imports

	sort.Slice(importItems, func(i, j int) bool {
		return importItems[i].Href < importItems[j].Href
	})

	// Does not handle pattern matching
	for _, item := range importItems {

		// Shouldn't be both but we can still handle the scenario
		if item.IncludeControls != nil {
			includeControls := *item.IncludeControls
			for _, includeControl := range includeControls {
				includes := *includeControl.WithIds
				sort.Slice(includes, func(i, j int) bool {
					return includes[i] < includes[j]
				})
				includeControl.WithIds = &includes
			}
		}

		if item.ExcludeControls != nil {
			excludeControls := *item.ExcludeControls
			for _, excludeControl := range excludeControls {
				exclude := *excludeControl.WithIds
				sort.Slice(exclude, func(i, j int) bool {
					return exclude[i] < exclude[j]
				})
				excludeControl.WithIds = &exclude
			}
		}

	}

	return
}

func (p Profile) HandleExisting(filepath string) error {
	exists, err := common.CheckFileExists(filepath)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Output File %s currently exist - cannot merge artifacts\n", filepath)
	} else {
		return nil
	}
}

// NewAssessmentResults creates a new assessment results object from the given data.
func NewProfile(data []byte) (Profile, error) {
	var profile Profile

	var oscalModels oscalTypes.OscalModels

	err := multiModelValidate(data)
	if err != nil {
		return profile, err
	}

	err = yaml.Unmarshal(data, &oscalModels)
	if err != nil {
		return profile, err
	}

	profile.Model = oscalModels.Profile
	profile.ModelType = "profile"

	return profile, nil
}

func GenerateProfile(source string, include []string, exclude []string) (profile Profile, err error) {

	var model oscalTypes.Profile

	// Single time used for all time related fields
	rfc3339Time := time.Now()

	// Always create a new UUID for the assessment results (for now)
	model.UUID = uuid.NewUUID()

	// Create metadata object with requires fields and a few extras
	// Where do we establish what `version` should be?
	model.Metadata = oscalTypes.Metadata{
		Title:        "Profile",
		Version:      "0.0.1",
		OscalVersion: OSCAL_VERSION,
		Remarks:      "Profile generated from Lula",
		Published:    &rfc3339Time,
		LastModified: rfc3339Time,
	}

	// Include would include the specified controls and exclude the rest
	// Exclude would exclude the specified controls and include the rest
	// Both doesn't make sense - TODO: Need to validate what OSCAL supports here
	includedControls := []oscalTypes.SelectControlById{
		oscalTypes.SelectControlById{
			WithIds: &include,
		},
	}

	excludedControls := []oscalTypes.SelectControlById{
		oscalTypes.SelectControlById{
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

	profile.Model = &model
	profile.ModelType = "profile"

	return profile, nil

}

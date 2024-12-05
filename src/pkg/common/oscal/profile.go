package oscal

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"gopkg.in/yaml.v3"

	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
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

func (p *Profile) HandleExisting(path string) error {
	exists, err := common.CheckFileExists(path)
	if err != nil {
		return err
	}
	if exists {
		path = filepath.Clean(path)
		existingFileBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}
		profile := NewProfile()
		err = profile.NewModel(existingFileBytes)
		if err != nil {
			return err
		}
		model, err := MergeProfileModels(profile.Model, p.Model)
		if err != nil {
			return err
		}
		p.Model = model
		return nil
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

func GenerateProfile(command string, source string, include []string, exclude []string, all bool) (*Profile, error) {

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

	if all {
		importItem.IncludeAll = &map[string]interface{}{}
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

func MergeProfileModels(original *oscalTypes.Profile, latest *oscalTypes.Profile) (*oscalTypes.Profile, error) {

	originalMap := make(map[string]oscalTypes.Import)
	latestMap := make(map[string]oscalTypes.Import)

	// Merge import items by href
	// we need to account for 1 -> N items in both models
	// do a simple merge by href now -> TODO: add intelligence for storing some source identifying information
	if len(original.Imports) == 0 {
		return original, fmt.Errorf("existing profile has no imports")
	}

	for _, item := range original.Imports {
		originalMap[item.Href] = item
	}

	if len(latest.Imports) == 0 {
		return original, fmt.Errorf("new profile has no imports ")
	}

	for _, item := range latest.Imports {
		latestMap[item.Href] = item
	}

	tempImports := make([]oscalTypes.Import, 0)

	for key, value := range latestMap {
		if _, ok := originalMap[key]; ok {
			// item exists - replace
			tempImports = append(tempImports, value)
			delete(originalMap, key)
		} else {
			// append the item
			tempImports = append(tempImports, value)
		}
	}

	for _, item := range originalMap {
		tempImports = append(tempImports, item)
	}

	// merge the back-matter resources
	if original.BackMatter != nil && latest.BackMatter != nil {
		original.BackMatter = &oscalTypes.BackMatter{
			Resources: mergeResources(original.BackMatter.Resources, latest.BackMatter.Resources),
		}
	} else if original.BackMatter == nil && latest.BackMatter != nil {
		original.BackMatter = latest.BackMatter
	}

	original.Imports = tempImports
	original.Metadata.LastModified = time.Now()

	// Artifact will be modified - need to update the UUID
	original.UUID = uuid.NewUUID()

	return original, nil
}

type ControlMap map[string]oscalTypes.Control

// ResolveProfileControls resolves all controls in the profile by checking any imported profiles or catalogs
// Returns a map[string]ControlMap where the key is the UUID of the source that dictates the controls (profile or catalog)
// NOTE: Profiles that contain Hrefs as references to internal identifiers (e.g., "#<UUID>") cannot currently be resolved
func ResolveProfileControls(profile *oscalTypes.Profile, profilePath, rootDir string, include, exclude []string) (map[string]ControlMap, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile is nil")
	}

	sourceControlMap := make(map[string]ControlMap)

	// Resolve the directory for imports
	importDir := network.GetLocalFileDir(profilePath, rootDir)

	for _, importItem := range profile.Imports {
		importedSourceControlMap, err := controlsFromImport(importItem, importDir)
		if err != nil {
			return nil, err
		}

		// Update sourceControlMap for profile and imports
		for source, controlMap := range importedSourceControlMap {
			addedControlMap := make(map[string]oscalTypes.Control)

			// Drop any excluded controls; include only included controls
			for id, control := range controlMap {
				if !AddControl(id, include, exclude) {
					continue
				}
				// If the source/control is not already in the map, add it
				if _, ok := addedControlMap[id]; !ok {
					addedControlMap[id] = control
				}
			}

			// Update sourceControlMap with anything imported
			sourceControlMap[source] = addedControlMap

			// Append controls to sourceControlMap[profile.UUID]
			if _, ok := sourceControlMap[profile.UUID]; !ok {
				sourceControlMap[profile.UUID] = make(map[string]oscalTypes.Control)
			}
			for id, control := range addedControlMap {
				if _, ok := sourceControlMap[profile.UUID][id]; !ok {
					sourceControlMap[profile.UUID][id] = control
				}
			}
		}
	}

	return sourceControlMap, nil
}

// Recursive function that resolves all controls in the provided profile and any imported profiles or catalogs
// rootDir is needed to resolve relative paths for imports
func controlsFromImport(importItem oscalTypes.Import, rootDir string) (controlMap map[string]ControlMap, err error) {
	// Fetch the import item
	var dataBytes []byte
	var fetchOpts []network.FetchOption
	if rootDir != "" {
		fetchOpts = append(fetchOpts, network.WithBaseDir(rootDir))
	}

	dataBytes, err = network.Fetch(importItem.Href, fetchOpts...)
	if err != nil {
		return controlMap, err
	}

	oscalModel, err := NewOscalModel(dataBytes)
	if err != nil {
		return controlMap, err
	}

	// Get include and exclude controls from import item
	var include []string
	var exclude []string
	if importItem.IncludeControls != nil {
		for _, includeControl := range *importItem.IncludeControls {
			include = append(include, *includeControl.WithIds...)
		}
	} else if importItem.ExcludeControls != nil {
		for _, excludeControl := range *importItem.ExcludeControls {
			exclude = append(exclude, *excludeControl.WithIds...)
		}
	}

	modelType, err := GetOscalModel(oscalModel)
	if err != nil {
		return controlMap, err
	}
	switch modelType {
	case "profile":
		return ResolveProfileControls(oscalModel.Profile, importItem.Href, rootDir, include, exclude)
	case "catalog":
		catalogControls, err := ResolveCatalogControls(oscalModel.Catalog, include, exclude)
		if err != nil {
			return nil, err
		}
		return map[string]ControlMap{
			oscalModel.Catalog.UUID: catalogControls,
		}, nil
	}

	return nil, nil
}

// AddControl takes the control-id, include and exclude lists and returns a boolean indicating if the control should be included
func AddControl(controlId string, include, exclude []string) bool {
	// If include is not empty, check if the control is included, default to AddControl=false
	// Include takes precedence over exclude
	if len(include) > 0 {
		for _, includeControl := range include {
			if controlId == includeControl {
				return true
			}
		}
		return false
	}

	// If exclude is not empty, check if the control is excluded, default to AddControl=true
	if len(exclude) > 0 {
		for _, excludeControl := range exclude {
			if controlId == excludeControl {
				return false
			}
		}
		return true
	}

	// If neither include nor exclude is specified, return AddControl=true
	return true
}

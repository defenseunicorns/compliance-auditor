package composition

import (
	"fmt"

	oscalTypes_1_1_2 "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/network"
)

// ResourceStore is a store of resources.
type ResourceStore struct {
	existing  map[string]*oscalTypes_1_1_2.Resource
	fetched   map[string]*oscalTypes_1_1_2.Resource
	hrefIdMap map[string][]string
}

// NewResourceStoreFromBackMatter creates a new resource store from the back matter of a component definition.
func NewResourceStoreFromBackMatter(backMatter *oscalTypes_1_1_2.BackMatter) *ResourceStore {
	store := &ResourceStore{
		existing: make(map[string]*oscalTypes_1_1_2.Resource),
		fetched:  make(map[string]*oscalTypes_1_1_2.Resource),
	}

	if backMatter != nil && *backMatter.Resources != nil {
		for _, resource := range *backMatter.Resources {
			store.AddExisting(&resource)
		}
	}

	return store
}

// AddExisting adds a resource to the store that is already in the back matter.
func (s *ResourceStore) AddExisting(resource *oscalTypes_1_1_2.Resource) {
	s.existing[resource.UUID] = resource
}

// GetExisting returns the resource with the given ID, if it exists.
func (s *ResourceStore) GetExisting(id string) (*oscalTypes_1_1_2.Resource, bool) {
	resource, ok := s.existing[id]
	return resource, ok
}

// AddFetched adds a resource to the store that was fetched from a remote source.
func (s *ResourceStore) AddFetched(resource *oscalTypes_1_1_2.Resource) {
	s.fetched[resource.UUID] = resource
}

// GetFetched returns the resource that was fetched from a remote source with the given ID, if it exists.
func (s *ResourceStore) GetFetched(id string) (*oscalTypes_1_1_2.Resource, bool) {
	resource, ok := s.fetched[id]
	return resource, ok
}

// AllFetched returns all the resources that were fetched from a remote source.
func (s *ResourceStore) AllFetched() []oscalTypes_1_1_2.Resource {
	resources := make([]oscalTypes_1_1_2.Resource, 0, len(s.fetched))
	for _, resource := range s.fetched {
		resources = append(resources, *resource)
	}
	return resources
}

// SetHrefIds sets the resource ids for a given href
func (s *ResourceStore) SetHrefIds(href string, ids []string) {
	s.hrefIdMap[href] = ids
}

// GetHrefIds gets the resource ids for a given href
func (s *ResourceStore) GetHrefIds(href string) (ids []string, err error) {
	if ids, ok := s.hrefIdMap[href]; ok {
		return ids, nil
	}
	return nil, fmt.Errorf("href #%s not found", href)
}

// Get returns the resource with the given ID, if it exists.
func (s *ResourceStore) Get(id string) (*oscalTypes_1_1_2.Resource, bool) {
	resource, inExisting := s.GetExisting(id)
	if inExisting {
		return resource, true
	}

	resource, inFetched := s.GetFetched(id)
	return resource, inFetched
}

// Has returns true if the resource store has a resource with the given ID.
func (s *ResourceStore) Has(id string) bool {
	_, inExisting := s.existing[id]
	_, inFetched := s.fetched[id]
	return inExisting || inFetched
}

// AddFromLink adds resources from a link to the store.
func (s *ResourceStore) AddFromLink(link *oscalTypes_1_1_2.Link) (ids []string, err error) {
	if link == nil {
		return nil, fmt.Errorf("link is nil")
	}
	id := common.TrimIdPrefix(link.Href)

	if link.ResourceFragment != common.WILDCARD && link.ResourceFragment != "" {
		id = common.TrimIdPrefix(link.ResourceFragment)
	}

	if s.Has(id) {
		return []string{id}, err
	}

	if ids, err = s.GetHrefIds(id); err == nil {
		return ids, err
	}

	return s.fetchFromRemoteLink(link)
}

func (s *ResourceStore) fetchFromRemoteLink(link *oscalTypes_1_1_2.Link) (ids []string, err error) {
	wantedId := common.TrimIdPrefix(link.ResourceFragment)

	validationBytes, err := network.Fetch(link.Href)
	if err != nil {
		return nil, err
	}

	validationArr, err := common.ReadValidationsFromYaml(validationBytes)
	if err != nil {
		return nil, err
	}
	isSingleValidation := len(validationArr) == 1

	for _, validation := range validationArr {
		resource, err := validation.ToResource()
		if err != nil {
			return nil, err
		}

		s.AddFetched(resource)

		if wantedId == resource.UUID || wantedId == common.WILDCARD || isSingleValidation {
			ids = append(ids, resource.UUID)
		}
	}

	return ids, err
}

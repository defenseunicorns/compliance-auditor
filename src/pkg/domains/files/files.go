package files

import (
	"fmt"

	"github.com/defenseunicorns/lula/src/types"
	"github.com/open-policy-agent/conftest/parser"
)

type Domain struct {
	Spec *Spec
}

// GetResources gathers the input files to be tested.
func (d Domain) GetResources() (types.DomainResources, error) {
	// conftest's parser returns a map[string]interface where the filenames are
	// the primary map keys.
	return parser.ParseConfigurations(d.Spec.Filepaths)
}

// IsExecutable returns false; the file domain is read-only.
func (d Domain) IsExecutable() bool { return false }

func CreateDomain(spec *Spec) (types.Domain, error) {
	if len(spec.Filepaths) == 0 {
		return nil, fmt.Errorf("file-spec must not be empty")
	}
	return Domain{spec}, nil
}

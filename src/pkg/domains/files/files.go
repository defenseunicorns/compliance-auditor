package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/types"
	"github.com/open-policy-agent/conftest/parser"
)

type Domain struct {
	Spec *Spec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// GetResources gathers the input files to be tested.
func (d Domain) GetResources() (types.DomainResources, error) {
	// see TODO below: maybe this is a REAL directory?
	dst, err := os.MkdirTemp("", "lula-files")
	if err != nil {
		return nil, err
	}

	// TODO? this might be a nice configurable option (for debugging) - store
	// the files into a local .lula directory that doesn't necessarily get
	// removed.
	defer os.RemoveAll(dst)

	// Copy files to a temporary location
	for _, path := range d.Spec.Filepaths {
		bytes, err := network.Fetch(path.Path)
		if err != nil {
			return nil, fmt.Errorf("error getting source files: %w", err)
		}
		os.WriteFile(filepath.Join(dst, path.Name), bytes, 0666)
	}

	// get a list of all the files we just downloaded in the temporary directory
	files := make([]string, 0)
	err = filepath.WalkDir(dst, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking downloaded file tree: %w", err)
	}

	// conftest's parser returns a map[string]interface where the filenames are
	// the primary map keys.
	config, err := parser.ParseConfigurations(files)
	if err != nil {
		return nil, err
	}

	// clean up the resources so it's just using the filename
	drs := make(types.DomainResources, len(config))
	for k, v := range config {
		rel, err := filepath.Rel(dst, k)
		if err != nil {
			return nil, fmt.Errorf("error determining relative file path: %w", err)
		}
		drs[rel] = v
	}
	return drs, nil
}

// IsExecutable returns false; the file domain is read-only.
//
// The files domain will download remote files into a temporary directory if the
// file paths are remote, but that is temporary and it is not mutating existing
// resources.
func (d Domain) IsExecutable() bool { return false }

func CreateDomain(spec *Spec) (types.Domain, error) {
	if len(spec.Filepaths) == 0 {
		return nil, fmt.Errorf("file-spec must not be empty")
	}
	return Domain{spec}, nil
}

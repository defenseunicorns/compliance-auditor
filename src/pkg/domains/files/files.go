package files

import (
	"context"
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
func (d Domain) GetResources(ctx context.Context) (types.DomainResources, error) {
	var workDir string
	var ok bool
	if workDir, ok = ctx.Value(types.LulaValidationWorkDir).(string); !ok {
		// if unset, assume lula is working in the same directory the inputFile is in
		workDir = "."
	}

	// see TODO below: maybe this is a REAL directory?
	dst, err := os.MkdirTemp("", "lula-files")
	if err != nil {
		return nil, err
	}

	// TODO? this might be a nice configurable option (for debugging) - store
	// the files into a local .lula directory that doesn't necessarily get
	// removed.
	defer os.RemoveAll(dst)

	// make a map of rel filepaths to the user-supplied name, so we can re-key the DomainResources later on.
	filenames := make(map[string]string, len(d.Spec.Filepaths))

	// Copy files to a temporary location
	for _, path := range d.Spec.Filepaths {
		file := filepath.Join(workDir, path.Path)
		bytes, err := network.Fetch(file)
		if err != nil {
			return nil, fmt.Errorf("error getting source files: %w", err)
		}

		// We'll just use the filename when writing the file so it's easier to reference later
		relname := filepath.Base(path.Path)

		err = os.WriteFile(filepath.Join(dst, relname), bytes, 0666)
		if err != nil {
			return nil, fmt.Errorf("error writing local files: %w", err)
		}
		// and save this info for later
		filenames[relname] = path.Name
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

	// clean up the resources so it's using the filepath.Name as the map key,
	// instead of the file path.
	drs := make(types.DomainResources, len(config))
	for k, v := range config {
		rel, err := filepath.Rel(dst, k)
		if err != nil {
			return nil, fmt.Errorf("error determining relative file path: %w", err)
		}
		drs[filenames[rel]] = v
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

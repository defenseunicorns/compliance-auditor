package files

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/lula/src/types"
	getter "github.com/hashicorp/go-getter/v2"
	"github.com/open-policy-agent/conftest/parser"
)

type Domain struct {
	Spec *Spec
}

// GetResources gathers the input files to be tested.
func (d Domain) GetResources() (types.DomainResources, error) {
	// see TODO below: maybe this is a REAL directory?
	dst, err := os.MkdirTemp("", "lula-files")
	if err != nil {
		return nil, err
	}

	// TODO? this might be a nice configurable option (for debugging) - or we go
	// the terraform route (sorry) and download and store them into a local
	// .lula directory
	defer os.RemoveAll(dst)

	// Copy file to a temporarly location, using go-getter to pull any remote files
	// TODO: use a real context, this isn't correct.
	g := getter.DefaultClient
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to determine current directory: %w", err)
	}

	for _, filepath := range d.Spec.Filepaths {
		_, err := g.Get(context.Background(), &getter.Request{
			Src: filepath,
			Dst: dst,
			Pwd: pwd,
		})
		if err != nil {
			return nil, fmt.Errorf("error getting source files: %w", err)
		}
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

	//clean up the resources so it's just using the filename
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
// This is a tiny lie: the files domain will download remote files into a
// temporary directory if the file paths are remote, but that is temporary and
// it is not mutating existing resources so we're calling it non-executable.
func (d Domain) IsExecutable() bool { return false }

func CreateDomain(spec *Spec) (types.Domain, error) {
	if len(spec.Filepaths) == 0 {
		return nil, fmt.Errorf("file-spec must not be empty")
	}
	return Domain{spec}, nil
}

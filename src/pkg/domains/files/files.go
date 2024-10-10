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
	filenames := make(map[string]string, 0)

	// unstructuredFiles is used to store a list of files that Lula needs to parse.
	unstructuredFiles := make([]FileInfo, 0)
	filesWithParsers := make(map[string][]FileInfo, 0)

	// Copy files to a temporary location
	for _, fi := range d.Spec.Filepaths {
		if fi.Parser != "" {
			if fi.Parser == "string" {
				unstructuredFiles = append(unstructuredFiles, fi)
				continue
			} else {
				filesWithParsers[fi.Parser] = append(filesWithParsers[fi.Parser], fi)
				continue
			}
		}

		file := filepath.Join(workDir, fi.Path)
		relname, err := copyFile(dst, file)
		if err != nil {
			return nil, fmt.Errorf("error writing local files: %w", err)
		}

		// and save this info for later
		filenames[relname] = fi.Name
	}

	// get a list of all the files we just downloaded in the temporary directory
	files, err := listFiles(dst)
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
	drs := make(types.DomainResources, len(config)+len(unstructuredFiles)+len(filesWithParsers))
	for k, v := range config {
		rel, err := filepath.Rel(dst, k)
		if err != nil {
			return nil, fmt.Errorf("error determining relative file path: %w", err)
		}
		drs[filenames[rel]] = v
	}

	// Now for the custom parsing: user-specified parsers and string files.

	for parserName, filesByParser := range filesWithParsers {
		// make a sub directory by parser name
		parserDir, err := os.MkdirTemp(dst, parserName)
		if err != nil {
			return nil, err
		}

		for _, fi := range filesByParser {
			file := filepath.Join(workDir, fi.Path)
			relname, err := copyFile(parserDir, file)
			if err != nil {
				return nil, fmt.Errorf("error writing local files: %w", err)
			}

			// and save this info for later
			filenames[relname] = fi.Name
		}

		// get a list of all the files we just downloaded in the temporary directory
		files, err := listFiles(parserDir)
		if err != nil {
			return nil, fmt.Errorf("error walking downloaded file tree: %w", err)
		}

		parsedConfig, err := parser.ParseConfigurationsAs(files, parserName)
		if err != nil {
			return nil, err
		}

		for k, v := range parsedConfig {
			rel, err := filepath.Rel(parserDir, k)
			if err != nil {
				return nil, fmt.Errorf("error determining relative file path: %w", err)
			}
			drs[filenames[rel]] = v
		}
	}

	// add the string form of the unstructured files
	for _, f := range unstructuredFiles {
		// we don't need to copy these files, we'll just slurp the contents into
		// a string and append that as one big DomainResource
		b, err := os.ReadFile(filepath.Join(workDir, f.Path))
		if err != nil {
			return nil, fmt.Errorf("error reading source files: %w", err)
		}
		drs[f.Name] = string(b)
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

// copyFile is a helper function that copies a file from source to dst, and returns the relative file path between the two.
func copyFile(dst string, src string) (string, error) {
	bytes, err := network.Fetch(src)
	if err != nil {
		return "", fmt.Errorf("error getting source files: %w", err)
	}

	// We'll use the filename when writing the file so it's easier to reference later
	relname := filepath.Base(src)

	return relname, os.WriteFile(filepath.Join(dst, relname), bytes, 0666)
}

func listFiles(dir string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

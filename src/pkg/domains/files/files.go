package files

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/open-policy-agent/conftest/parser"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/types"
)

type Domain struct {
	Spec *Spec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// GetResources gathers the input files to be tested.
func (d Domain) GetResources(ctx context.Context) (types.DomainResources, error) {
	var workDir string
	var ok bool
	var errs error
	tmpDRs := make(map[string]interface{})
	if workDir, ok = ctx.Value(types.LulaValidationWorkDir).(string); !ok {
		// if unset, assume lula is already working in the same directory the inputFile is in
		workDir = "."
	}

	dst, err := os.MkdirTemp("", "lula-files-")
	if err != nil {
		// allow returning on error here?
		return nil, err
	}

	// (potential) TODO: this could be a nice configurable option (for
	// debugging) - store the files into a local .lula directory that doesn't
	// get removed.
	defer os.RemoveAll(dst)

	// filenames stores a map of relative filenames to the user-supplied Name,
	// so we can re-key the DomainResources later on.
	filenames := make(map[string]string, 0)
	// unstructuredFiles is used to store a list of files that Lula needs to
	// parse as strings.
	unstructuredFiles := make([]FileInfo, 0)
	// filesWithParsers stores files with user-specified parsers to pass to
	// conftest.
	filesWithParsers := make(map[string][]FileInfo, 0)

	// Copy files to a temporary location. In this loop we only grab files that
	// don't have configured parsers.
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

		realdst := filepath.Join(dst, filepath.Base(fi.Path))
		filename, err := copyFile(ctx, realdst, fi.Path, workDir)
		if err != nil {
			// Assign empty data value for reporting purposes
			tmpDRs[fi.Name] = map[string]interface{}{}
			errs = errors.Join(errs, fmt.Errorf("error writing local files: %w", err))
			continue
		}

		// and save this info for later
		filenames[filename] = fi.Name
	}

	// get a list of all the files we just downloaded in the temporary directory
	files, err := listFiles(dst)
	if err != nil {
		return tmpDRs, errors.Join(errs, fmt.Errorf("error walking downloaded file tree: %w", err))
	}

	// conftest's parser returns a map[string]interface where the filenames are
	// the primary map keys.
	// need to test this to understand the outcomes on a single file error on the return values
	config, err := parser.ParseConfigurations(files)
	// Copy values over to the temporary domain resources
	for k, v := range config {
		tmpDRs[k] = v
	}
	if err != nil {
		errs = errors.Join(errs, err)
		return tmpDRs, errs
	}

	// clean up the resources so it's using the filepath.Name as the map key,
	// instead of the file path.
	drs := make(types.DomainResources, len(tmpDRs)+len(unstructuredFiles)+len(filesWithParsers))
	for k, v := range tmpDRs {
		rel, err := filepath.Rel(dst, k)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("error determining relative file path: %w", err))
			drs[k] = v
			continue
		}
		drs[filenames[rel]] = v
	}

	// Now for the custom parsing: user-specified parsers and string files.
	for parserName, filesByParser := range filesWithParsers {
		// make a sub directory by parser name
		parserDir, err := os.MkdirTemp(dst, parserName)
		if err != nil {
			return drs, err
		}

		for _, fi := range filesByParser {
			dst := filepath.Join(parserDir, filepath.Base(fi.Path))
			relname, err := copyFile(ctx, dst, fi.Path, workDir)
			if err != nil {
				drs[fi.Name] = map[string]interface{}{}
				errs = errors.Join(errs, fmt.Errorf("error writing local files: %w", err))
			}

			// and save this info for later
			filenames[relname] = fi.Name
		}

		// get a list of all the files we just downloaded in the temporary directory
		files, err := listFiles(parserDir)
		if err != nil {
			return drs, errors.Join(errs, fmt.Errorf("error walking downloaded file tree: %w", err))
		}

		parsedConfig, err := parser.ParseConfigurationsAs(files, parserName)
		if err != nil {
			return drs, err
		}

		for k, v := range parsedConfig {
			rel, err := filepath.Rel(parserDir, k)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("error determining relative file path: %w", err))
				drs[filenames[k]] = v
				continue
			}
			drs[filenames[rel]] = v
		}
	}

	// add the string form of the unstructured files
	for _, f := range unstructuredFiles {
		stringdir, err := os.MkdirTemp(dst, "string")
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("error reading source files: %w", err))
			drs[f.Name] = ""
			continue
		}

		dst := filepath.Clean(filepath.Join(stringdir, f.Path))
		_, err = copyFile(ctx, dst, f.Path, workDir)
		if err != nil {
			return nil, fmt.Errorf("error writing local files: %w", err)
		}

		b, err := os.ReadFile(dst)
		if err != nil {
			return nil, fmt.Errorf("error reading local file: %w", err)
		}

		drs[f.Name] = string(b)
	}

	return drs, errs
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

// copyFile is a helper function that copies a file from source to dst, and
// returns the base filename.
func copyFile(ctx context.Context, dst, src, wd string) (string, error) {
	realdst, err := network.DownloadFile(ctx, dst, src, wd)
	if err != nil {
		return "", fmt.Errorf("error getting source files: %w", err)
	}
	return filepath.Base(realdst), nil
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

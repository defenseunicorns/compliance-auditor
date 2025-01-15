package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	goversion "github.com/hashicorp/go-version"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	"github.com/defenseunicorns/lula/src/pkg/domains/files"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
)

const (
	UUID_PREFIX    = "#"
	WILDCARD       = "*"
	YAML_DELIMITER = "---"
)

// TrimIdPrefix trims the id prefix from the given id
func TrimIdPrefix(id string) string {
	return strings.TrimPrefix(id, UUID_PREFIX)
}

// AddIdPrefix adds the id prefix to the given id
func AddIdPrefix(id string) string {
	return fmt.Sprintf("%s%s", UUID_PREFIX, id)
}

// IsLulaLink checks if the link is a lula link
func IsLulaLink(link oscalTypes.Link) bool {
	rel := strings.Split(link.Rel, ".")
	return link.Text == "Lula Validation" || rel[0] == "lula"
}

// ReadFileToBytes reads a file to bytes
func ReadFileToBytes(path string) ([]byte, error) {
	var data []byte
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return data, fmt.Errorf("Path: %v does not exist - unable to digest document", path)
	}
	path = filepath.Clean(path)
	data, err = os.ReadFile(path)
	if err != nil {
		return data, err
	}

	return data, nil
}

// ReadValidationsFromYaml reads a yaml file of validations to an array of validations
func ReadValidationsFromYaml(validationBytes []byte) (validations []Validation, err error) {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(validationBytes), 4096)

	for {
		validation := &Validation{}
		if err := decoder.Decode(validation); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		validations = append(validations, *validation)
	}

	return validations, nil
}

// Returns version validity
func IsVersionValid(versionConstraint string, version string) (bool, error) {
	if version == "unset" {
		// Default cli version is "unset", enabling users to run directly from source code
		// This is not a valid version, but we want to allow it for development purposes
		return true, nil
	}

	currentVersion, err := goversion.NewVersion(version)
	if err != nil {
		return false, err
	}
	constraints, err := goversion.NewConstraint(versionConstraint)
	if err != nil {
		return false, err
	}
	if constraints.Check(currentVersion) {
		return true, nil
	}
	return false, nil
}

// SetCwdToFileDir takes a path and changes the current working directory to the directory of the path
func SetCwdToFileDir(dirPath string) (resetFunc func(), err error) {
	// Bail if empty
	if dirPath == "" {
		return resetFunc, fmt.Errorf("dirPath is empty")
	}
	info, err := os.Stat(dirPath)
	// Bail if the path does not exist
	if err != nil {
		return resetFunc, err
	}
	// Get the directory of the path if it is not a directory
	if !info.IsDir() {
		dirPath = filepath.Dir(dirPath)
	}

	// Save the current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		return resetFunc, err
	}

	// Change the current working directory
	err = os.Chdir(dirPath)
	if err != nil {
		return resetFunc, err
	}

	// Return a function to change the current working directory back to the original directory
	resetFunc = func() {
		err = os.Chdir(originalDir)
		if err != nil {
			message.Warnf("unable to change cwd back to %s: %v", originalDir, err)
		}
	}

	// Change back to the original working directory when done
	return resetFunc, err
}

// Get the domain and providers
func GetDomain(domain *Domain) (types.Domain, error) {
	if domain == nil {
		return nil, fmt.Errorf("domain is nil")
	}
	switch domain.Type {
	case "kubernetes":
		return kube.CreateKubernetesDomain(domain.KubernetesSpec)
	case "api":
		return api.CreateApiDomain(domain.ApiSpec)
	case "file":
		return files.CreateDomain(domain.FileSpec)
	default:
		return nil, fmt.Errorf("domain is unsupported")
	}
}

func GetProvider(provider *Provider, ctx context.Context) (types.Provider, error) {
	if provider == nil {
		return nil, fmt.Errorf("provider is nil")
	}
	switch provider.Type {
	case "opa":
		return opa.CreateOpaProvider(ctx, provider.OpaSpec)
	case "kyverno":
		return kyverno.CreateKyvernoProvider(ctx, provider.KyvernoSpec)
	default:
		return nil, fmt.Errorf("provider is unsupported")
	}
}

// Converts a raw string to a Validation object (string -> common.Validation -> types.Validation)
func ValidationFromString(raw, uuid string) (validation types.LulaValidation, err error) {
	if raw == "" {
		return validation, fmt.Errorf("validation string is empty")
	}

	var validationData Validation
	err = validationData.UnmarshalYaml([]byte(raw))
	if err != nil {
		return validation, err
	}

	validation, err = validationData.ToLulaValidation(uuid)
	if err != nil {
		return validation, err
	}

	return validation, nil
}

func CheckFileExists(filepath string) (bool, error) {
	if _, err := os.Stat(filepath); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

// CleanMultilineString removes leading and trailing whitespace from a multiline string
func CleanMultilineString(str string) string {
	re := regexp.MustCompile(`[ \t]+\r?\n`)
	formatted := re.ReplaceAllString(str, "\n")
	return formatted
}

// RemapPath takes an input path, relative to the baseDir, and remaps it to be relative to the newDir
// Example: path = "folder/file.txt", baseDir = "/home/user/dir", newDir = "/home/user/newDir"
// output path = "../dir/folder/file.txt"
func RemapPath(path string, baseDir string, newDir string) (string, error) {
	// Do nothing if the path is a UUID reference
	if isUUIDReference(path) {
		return path, nil
	}

	// Return if the path is a URL or absolute link
	localDir := network.GetLocalFileDir(path, baseDir)
	if localDir == "" {
		return path, nil
	}

	// Trim file://, if present
	path = strings.TrimPrefix(path, "file://")

	// Find the relative path from newDir to baseDir
	relativePath, err := filepath.Rel(newDir, baseDir)
	if err != nil {
		return "", err
	}

	// Append the original relative path to the computed relative path
	remappedPath := filepath.Join(relativePath, path)
	remappedPath = filepath.Clean(remappedPath)

	return remappedPath, nil
}

func isUUIDReference(path string) bool {
	path = strings.TrimPrefix(path, UUID_PREFIX)
	return checkValidUuid(path)
}

// TraverseAndUpdatePaths uses reflection to traverse the obj based on the path and update file path references
func TraverseAndUpdatePaths(obj interface{}, path string, baseDir string, newDir string) error {
	// Split the path into components
	components := splitPath(path)

	// Start reflection traversal
	return reflectTraverseAndUpdate(reflect.ValueOf(obj), components, baseDir, newDir)
}

func reflectTraverseAndUpdate(val reflect.Value, components []string, baseDir string, newDir string) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer
	}

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		// Handle slices/arrays
		for i := 0; i < val.Len(); i++ {
			err := reflectTraverseAndUpdate(val.Index(i), components, baseDir, newDir)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if val.Kind() == reflect.Struct && len(components) > 0 {
		// Handle structs
		field := val.FieldByName(components[0])
		if !field.IsValid() {
			return fmt.Errorf("field %s not found", components[0])
		}
		return reflectTraverseAndUpdate(field, components[1:], baseDir, newDir)
	}

	if len(components) == 0 {
		if val.Kind() == reflect.String {
			// Update the final field (assumed to be a string)
			newValue, err := RemapPath(val.String(), baseDir, newDir)
			if err != nil {
				return fmt.Errorf("error remapping path %s: %v", val.String(), err)
			}
			if val.CanSet() {
				val.SetString(newValue)
				return nil
			}
			return fmt.Errorf("unable to set string value")
		}
		// if val.Kind() is not a string, we can't update it
		return fmt.Errorf("cannot update type %s", val.Kind())
	}

	return nil
}

func splitPath(path string) []string {
	components := strings.Split(path, ".")
	return components
}

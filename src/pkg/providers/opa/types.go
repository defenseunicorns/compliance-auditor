package opa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/defenseunicorns/lula/src/pkg/common/network"
	"github.com/defenseunicorns/lula/src/types"
)

var (
	ErrNilSpec                = fmt.Errorf("spec is nil")
	ErrEmptyRego              = fmt.Errorf("rego policy cannot be empty")
	ErrInvalidValidationPath  = fmt.Errorf("validation field must be a json path")
	ErrInvalidObservationPath = fmt.Errorf("observation field must be a json path")
	ErrDownloadModule         = fmt.Errorf("error downloading module")
	ErrReadModule             = fmt.Errorf("error reading module")
)

type OpaProvider struct {
	// Spec is the specification of the OPA policy
	Spec *OpaSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

func CreateOpaProvider(_ context.Context, spec *OpaSpec) (types.Provider, error) {
	// Check validity of spec
	if spec == nil {
		return nil, ErrNilSpec
	}

	if spec.Rego == "" {
		return nil, ErrEmptyRego
	}

	if spec.Output != nil {
		if spec.Output.Validation != "" {
			if !strings.Contains(spec.Output.Validation, ".") {
				return nil, ErrInvalidValidationPath
			}
		}
		if spec.Output.Observations != nil {
			for _, observation := range spec.Output.Observations {
				if !strings.Contains(observation, ".") {
					return nil, ErrInvalidObservationPath
				}
			}
		}
	}

	return OpaProvider{
		Spec: spec,
	}, nil
}

// loadModules downloads the modules specified in the modulePaths map and returns
// a map of the module name to the module content.
func loadModules(ctx context.Context, modulePaths map[string]string) (map[string]string, error) {
	if len(modulePaths) == 0 {
		return nil, nil
	}

	workDir, ok := ctx.Value(types.LulaValidationWorkDir).(string)
	if !ok { // if unset, assume lula is already working in the same directory the inputFile is in
		workDir = "."
	}

	dst, err := os.MkdirTemp("", "lula-modules-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dst)

	loadedModules := make(map[string]string)
	for name, src := range modulePaths {
		dst := filepath.Join(dst, filepath.Base(src))
		tmp, err := network.DownloadFile(ctx, dst, src, workDir)
		if err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrDownloadModule, name, err)
		}
		content, err := os.ReadFile(filepath.Clean(tmp))
		if err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrReadModule, name, err)
		}
		loadedModules[name] = string(content)
	}

	return loadedModules, nil
}

func (o OpaProvider) Evaluate(ctx context.Context, resources types.DomainResources) (types.Result, error) {
	modules, err := loadModules(ctx, o.Spec.Modules)
	if err != nil {
		return types.Result{}, err
	}
	results, err := GetValidatedAssets(ctx, o.Spec.Rego, modules, resources, o.Spec.Output)
	if err != nil {
		return types.Result{}, err
	}
	return results, nil
}

// OpaSpec is the specification of the OPA policy, required if the provider type is opa
type OpaSpec struct {
	// Required: Rego is the OPA policy
	Rego string `json:"rego" yaml:"rego"`
	// Optional: Modules is a map of additional OPA modules to include. The key is the name of the
	// module and the value is the file with the contents of the module.
	Modules map[string]string `json:"modules,omitempty" yaml:"modules,omitempty"`
	// Optional: Output is the output of the OPA policy
	Output *OpaOutput `json:"output,omitempty" yaml:"output,omitempty"`
}

// OpaOutput Defines the output structure for OPA validation results, including validation status and additional observations.
type OpaOutput struct {
	// optional: Specifies the JSON path to a boolean value indicating the validation result.
	Validation string `json:"validation" yaml:"validation"`
	// optional: any additional observations to include (fields must resolve to strings)
	Observations []string `json:"observations" yaml:"observations"`
}

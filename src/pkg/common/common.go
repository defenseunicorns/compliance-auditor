package common

import (
	"context"
	"fmt"
	"os"

	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
	"github.com/defenseunicorns/lula/src/types"
	goversion "github.com/hashicorp/go-version"
)

func ReadFileToBytes(path string) ([]byte, error) {
	var data []byte
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return data, fmt.Errorf("Path: %v does not exist - unable to digest document", path)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		return data, err
	}

	return data, nil
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

// Get the domain and providers
func GetProvider(provider string, ctx context.Context) types.Provider {
	switch provider {
	case "opa":
		return opa.OpaProvider{
			Context: ctx,
		}
	case "kyverno":
		return kyverno.KyvernoProvider{
			Context: ctx,
		}
	default:
		return nil
	}
}

func GetDomain(domain string, ctx context.Context) types.Domain {
	switch domain {
	case "kubernetes":
		return kube.KubernetesDomain{
			Context: ctx,
		}
	case "api":
		return api.ApiDomain{}
	default:
		return nil
	}
}

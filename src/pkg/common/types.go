package common

import (
	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
)

// Data structures for ingesting validation data
type Validation struct {
	LulaVersion string   `json:"lula-version" yaml:"lula-version"`
	Metadata    Metadata `json:"metadata" yaml:"metadata"`
	Provider    Provider `json:"provider" yaml:"provider"`
	Domain      Domain   `json:"domain" yaml:"domain"`
}

// TODO: Perhaps extend this structure with other needed information, such as UUID or type of validation if workflow is needed
type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type Domain struct {
	Type           string              `json:"type" yaml:"type"`
	KubernetesSpec kube.KubernetesSpec `json:"kubernetes-spec" yaml:"kubernetes-spec"`
	ApiSpec        api.ApiSpec         `json:"api-spec" yaml:"api-spec"`
}

type Provider struct {
	Type        string              `json:"type" yaml:"type"`
	OpaSpec     opa.OpaSpec         `json:"opa-spec" yaml:"opa-spec"`
	KyvernoSpec kyverno.KyvernoSpec `json:"kyverno-spec" yaml:"kyverno-spec"`
}

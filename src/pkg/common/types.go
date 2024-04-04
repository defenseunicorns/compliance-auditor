package common

import (
	"github.com/defenseunicorns/lula/src/pkg/domains/api"
	kube "github.com/defenseunicorns/lula/src/pkg/domains/kubernetes"
	"github.com/defenseunicorns/lula/src/pkg/providers/kyverno"
	"github.com/defenseunicorns/lula/src/pkg/providers/opa"
)

// YAML data structures for ingesting validation data
type ValidationYaml struct {
	LulaVersion string       `json:"lula-version" yaml:"lula-version"`
	Metadata    MetadataYaml `json:"metadata" yaml:"metadata"`
	Target      TargetYaml   `json:"target" yaml:"target"`
}

type MetadataYaml struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
}

type TargetYaml struct {
	Provider ProviderYaml `json:"provider" yaml:"provider"`
	Domain   DomainYaml   `json:"domain" yaml:"domain"`
}

type DomainYaml struct {
	Type           string              `json:"type" yaml:"type"`
	KubernetesSpec kube.KubernetesSpec `json:"kubernetes-spec" yaml:"kubernetes-spec"`
	ApiSpec        api.ApiSpec         `json:"api-spec" yaml:"api-spec"`
}

type ProviderYaml struct {
	Type        string              `json:"type" yaml:"type"`
	OpaSpec     opa.OpaSpec         `json:"opa-spec" yaml:"opa-spec"`
	KyvernoSpec kyverno.KyvernoSpec `json:"kyverno-spec" yaml:"kyverno-spec"`
}

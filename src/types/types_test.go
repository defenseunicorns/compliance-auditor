package types_test

import (
	"testing"

	"github.com/defenseunicorns/lula/src/types"
)

func Test_Validation_Lint(t *testing.T) {

	tests := []struct {
		name       string
		validation types.Validation
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			validation: types.Validation{
				Title:       "Lula Validation",
				LulaVersion: ">= v0.1.0",
				Target: types.Target{
					Provider: "opa",
					Domain:   "kubernetes",
					Payload: types.Payload{
						Resources: []types.Resource{
							{
								Name: "podsvt",
								ResourceRule: types.ResourceRule{
									Group:      "core",
									Version:    "v1",
									Resource:   "pods",
									Namespaces: []string{"validation-test"},
								},
							},
						},
						Rego: "package validate\n\nimport future.keywords.every\n\nvalidate {\n  every pod in input.podsvt {\n  podLabel := pod.metadata.labels.foo\npodlabel == \"bar\"\n  }\n}",
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "error no validation",
			validation: types.Validation{},
			wantErr:    true,
		},
		{
			name: "error no validation.title",
			validation: types.Validation{
				Title: "Lula Validation",
			},
			wantErr: true,
		},
		{
			name: "error no target",
			validation: types.Validation{
				Title:       "Lula Validation",
				LulaVersion: ">= v0.1.0",
				Target:      types.Target{},
			},
			wantErr: true,
		},
		{
			name: "error resource-rule.name no resource-rule.namespaces",
			validation: types.Validation{
				Title:       "Lula Validation",
				LulaVersion: ">= v0.1.0",
				Target: types.Target{
					Provider: "opa",
					Domain:   "kubernetes",
					Payload: types.Payload{
						Resources: []types.Resource{
							{
								Name: "podsvt",
								ResourceRule: types.ResourceRule{
									Name:     "podsvt",
									Group:    "core",
									Version:  "v1",
									Resource: "pods",
								},
							},
						},
						Rego: "package validate\n\nimport future.keywords.every\n\nvalidate {\n  every pod in input.podsvt {\n  podLabel := pod.metadata.labels.foo\npodlabel == \"bar\"\n  }\n}",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error resource-rule.field no resource-rule.name",
			validation: types.Validation{
				Title:       "Lula Validation",
				LulaVersion: ">= v0.1.0",
				Target: types.Target{
					Provider: "opa",
					Domain:   "kubernetes",
					Payload: types.Payload{
						Resources: []types.Resource{
							{
								Name: "podsvt",
								ResourceRule: types.ResourceRule{
									Group:    "core",
									Version:  "v1",
									Resource: "pods",
									Field: types.Field{
										Jsonpath: "metadata.labels.foo",
										Type:     "json",
										Base64:   false,
									},
								},
							},
						},
						Rego: "package validate\n\nimport future.keywords.every\n\nvalidate {\n  every pod in input.podsvt {\n  podLabel := pod.metadata.labels.foo\npodlabel == \"bar\"\n  }\n}",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error no resource-rule.resource",
			validation: types.Validation{
				Title:       "Lula Validation",
				LulaVersion: ">= v0.1.0",
				Target: types.Target{
					Provider: "opa",
					Domain:   "kubernetes",
					Payload: types.Payload{
						Resources: []types.Resource{
							{
								Name: "podsvt",
								ResourceRule: types.ResourceRule{
									Group:   "core",
									Version: "v1",
								},
							},
						},
						Rego: "package validate\n\nimport future.keywords.every\n\nvalidate {\n  every pod in input.podsvt {\n  podLabel := pod.metadata.labels.foo\npodlabel == \"bar\"\n  }\n}",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error resource-rule no version",
			validation: types.Validation{
				Title:       "Lula Validation",
				LulaVersion: ">= v0.1.0",
				Target: types.Target{
					Provider: "opa",
					Domain:   "kubernetes",
					Payload: types.Payload{
						Resources: []types.Resource{
							{
								Name: "podsvt",
								ResourceRule: types.ResourceRule{
									Group:    "core",
									Resource: "pods",
								},
							},
						},
						Rego: "package validate\n\nimport future.keywords.every\n\nvalidate {\n  every pod in input.podsvt {\n  podLabel := pod.metadata.labels.foo\npodlabel == \"bar\"\n  }\n}",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.validation.Lint(); (err != nil) != tt.wantErr {
				t.Errorf("lintValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

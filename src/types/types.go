package types

type Validation struct {
	Title       string                 `json:"title" yaml:"title"`
	Description map[string]interface{} `json:"description" yaml:"description"`
	Evaluated   bool                   `json:"evaluated" yaml:"evaluated"`
	Result      Result                 `json:"result" yaml:"result"`
}

// native type for conversion to targeted report format
type Result struct {
	UUID        string `json:"uuid" yaml:"uuid"`
	ControlId   string `json:"control-id" yaml:"control-id"`
	Description string `json:"description" yaml:"description"`
	Passing     int    `json:"passing" yaml:"passing"`
	Failing     int    `json:"failing" yaml:"failing"`
	State       string `json:"state" yaml:"state"`
}

// Current placeholder for all requisite data in the payload
// Fields will be populated as required otherwise left empty
// This could be expanded as providers add more fields
type Payload struct {
	Resources []Resource `json:"resources" yaml:"resources"`
	Wait      Wait       `json:"wait" yaml:"wait"`
	Rego      string     `json:"rego" yaml:"rego"`
}

type Resource struct {
	Name         string       `json:"name" yaml:"name"`
	Description  string       `json:"description" yaml:"description"`
	ResourceRule ResourceRule `json:"resource-rule" yaml:"resource-rule"`
}

type Wait struct {
	Condition string `json:"condition" yaml:"condition"`
	Jsonpath  string `json:"jsonpath" yaml:"jsonpath"`
	Kind      string `json:"kind" yaml:"kind"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Timeout   string `json:"timeout" yaml:"timeout"`
}

type PayloadAPI struct {
	Requests []Request `mapstructure:"requests" json:"requests" yaml:"requests"`
	Rego     string    `json:"rego" yaml:"rego"`
}

type Request struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

type Target struct {
	Provider string  `json:"provider" yaml:"provider"`
	Domain   string  `json:"domain" yaml:"domain"`
	Payload  Payload `json:"payload" yaml:"payload"`
}

type ResourceRule struct {
	Name       string   `json:"name" yaml:"name"`
	Group      string   `json:"group" yaml:"group"`
	Version    string   `json:"version" yaml:"version"`
	Resource   string   `json:"resource" yaml:"resource"`
	Namespaces []string `json:"namespaces" yaml:"namespaces"`
}

// generate manifest types
type Manifest struct {
	Kind                string              `json:"kind" yaml:"kind"`
	Metadata            Metadata            `json:"metadata" yaml:"metadata"`
	ComponentDefinition ComponentDefinition `json:"component-definition" yaml:"component-definition"`
	SSP                 SSP                 `json:"system-security-plan" yaml:"system-security-plan"`
	SAP                 SAP                 `json:"system-assessment-plan" yaml:"system-assessment-plan"`
	POAM                POAM                `json:"plan-of-action-and-milestones" yaml:"plan-of-action-and-milestones"`
}

type Metadata struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

type ComponentDefinition struct {
	Name       string     `json:"name" yaml:"name"`
	Catalogs   []Catalog  `json:"catalogs" yaml:"catalogs"`
	Components []Artifact `json:"components" yaml:"components"`
}

type Catalog struct {
	Path     string    `json:"path" yaml:"path"`
	Url      string    `json:"url" yaml:"url"`
	Git      string    `json:"git" yaml:"git"`
	GitPath  string    `json:"gitPath" yaml:"gitPath"`
	Controls []Control `json:"controls" yaml:"controls"`
}

type Control struct {
	Id      string `json:"id" yaml:"id"`
	Remarks string `json:"remarks" yaml:"remarks"`
}

type SSP struct {
	Components []Artifact `json:"components" yaml:"components"`
}

type SAP struct {
	AssessmentResults []Artifact `json:"assessment-results" yaml:"assessment-results"`
}

type POAM struct {
	AssessmentResults []Artifact `json:"assessment-results" yaml:"assessment-results"`
}

type Artifact struct {
	Path    string `json:"path" yaml:"path"`
	Url     string `json:"url" yaml:"url"`
	Git     string `json:"git" yaml:"git"`
	GitPath string `json:"gitPath" yaml:"gitPath"`
}

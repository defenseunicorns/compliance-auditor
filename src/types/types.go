package types

// this is a temporary struct until a reporting model is selected
type ComplianceReport struct {
	UUID        string `json:"uuid" yaml:"uuid"`
	ControlId   string `json:"control-id" yaml:"control-id"`
	Description string `json:"description" yaml:"description"`
	Result      string `json:"result" yaml:"result"`
}

// The ReportObject keeps track of all pertinent information as it relates to relational data IE UUID's
// This will hopefully make transformation to the reporting model easier
// or be replaced by an OSCAL native type
type ReportObject struct {
	FilePaths   []string              `json:"file-paths" yaml:"file-paths"`
	UUID        string                `json:"uuid" yaml:"uuid"`
	Components  []Component           `json:"components" yaml:"components"`
	Validations map[string]Validation `json:"validations" yaml:"validations"`
}

type Validation struct {
	Title       string                 `json:"title" yaml:"title"`
	Description map[string]interface{} `json:"description" yaml:"description"`
	Evaluated   bool                   `json:"evaluated" yaml:"evaluated"`
	Result      Result                 `json:"result" yaml:"result"`
}

type Component struct {
	UUID                   string                  `json:"uuid" yaml:"uuid"`
	ControlImplementations []ControlImplementation `json:"control-implementations" yaml:"control-implementations"`
}

type ControlImplementation struct {
	UUID            string           `json:"uuid" yaml:"uuid"`
	ImplementedReqs []ImplementedReq `json:"implemented-reqs" yaml:"implemented-reqs"`
}

type ImplementedReq struct {
	UUID        string   `json:"uuid" yaml:"uuid"`
	ControlId   string   `json:"control-id" yaml:"control-id"`
	Status      string   `json:"status" yaml:"status"`
	Description string   `json:"description" yaml:"description"`
	Results     []Result `json:"results" yaml:"results"`
}

// native type for conversion to targeted report format
type Result struct {
	UUID        string `json:"uuid" yaml:"uuid"`
	ControlId   string `json:"control-id" yaml:"control-id"`
	Description string `json:"description" yaml:"description"`
	Passing     int    `json:"passing" yaml:"passing"`
	Failing     int    `json:"failing" yaml:"failing"`
	Result      string `json:"result" yaml:"result"`
}

// Current placeholder for all requisite data in the payload
// Fields will be populated as required otherwise left empty
// This could be expanded as providers add more fields
type Payload struct {
	ResourceRules []ResourceRule `json:"resource-rules" yaml:"resource-rules"`
	Rego          string         `json:"rego" yaml:"rego"`
}

type Target struct {
	Provider string  `json:"provider" yaml:"provider"`
	Domain   string  `json:"domain" yaml:"domain"`
	Payload  Payload `json:"payload" yaml:"payload"`
}

type ResourceRule struct {
	Group      string   `json:"group" yaml:"group"`
	Version    string   `json:"version" yaml:"version"`
	Resource   string   `json:"resource" yaml:"resource"`
	Namespaces []string `json:"namespaces" yaml:"namespaces"`
}

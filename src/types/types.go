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
	Rego      string     `json:"rego" yaml:"rego"`
}

type Resource struct {
	Name         string       `json:"name" yaml:"name"`
	Description  string       `json:"description" yaml:"description"`
	ResourceRule ResourceRule `json:"resource-rule" yaml:"resource-rule"`
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

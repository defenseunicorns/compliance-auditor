package template

type VariableConfig struct {
	Key       string
	Default   string
	Sensitive bool
}

type TemplateData struct {
	Constants          map[string]interface{}
	Variables          map[string]string
	SensitiveVariables map[string]string
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Constants:          make(map[string]interface{}),
		Variables:          make(map[string]string),
		SensitiveVariables: make(map[string]string),
	}
}

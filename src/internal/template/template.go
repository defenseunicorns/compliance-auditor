package template

import (
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	PREFIX = "LULA_VAR_"
	CONST  = "const"
	VAR    = "var"
)

func createTemplate() *template.Template {
	// Register custom template functions
	funcMap := template.FuncMap{
		"concatToRegoList": func(a []any) string {
			return concatToRegoList(a)
		},
		"mask": func(a string) string {
			return "********"
		},
		// Add more custom functions as needed
	}

	// Parse the template and apply the function map
	tpl := template.New("template").Funcs(funcMap)
	tpl.Option("missingkey=zero")
	return tpl
}

// ExecuteFullTemplate templates everything
func ExecuteFullTemplate(templateData *TemplateData, templateString string) ([]byte, error) {
	tpl := createTemplate()
	tpl, err := tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	allVars := MergeStringMaps(templateData.Variables, templateData.SensitiveVariables)
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: templateData.Constants,
		VAR:   allVars})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// ExecuteConstTemplate templates only constants
// this templates only values in the constants map
func ExecuteConstTemplate(constants map[string]interface{}, templateString string) ([]byte, error) {
	// Find anything {{ var.KEY }} and replace with {{ "{{ var.KEY }}" }}
	re := regexp.MustCompile(`{{\s*\.` + VAR + `\.([a-zA-Z0-9_]+)\s*}}`)
	templateString = re.ReplaceAllString(templateString, "{{ \"{{ ."+VAR+".$1 }}\" }}")

	tpl := createTemplate()
	tpl, err := tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: constants})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// ExecuteNonSensitiveTemplate templates only constants and non-sensitive variables
// used for compose operations
func ExecuteNonSensitiveTemplate(templateData *TemplateData, templateString string) ([]byte, error) {
	// Find any sensitive keys {{ var.KEY }}, where KEY is in templateData.SensitiveVariables and replace with {{ "{{ var.KEY }}" }}
	re := regexp.MustCompile(`{{\s*\.` + VAR + `\.([a-zA-Z0-9_]+)\s*}}`)
	varMatches := re.FindStringSubmatch(templateString)
	for _, m := range varMatches {
		if _, ok := templateData.SensitiveVariables[m]; ok {
			reSensitive := regexp.MustCompile(`{{\s*\.` + VAR + `\.` + m + `\s*}}`)
			templateString = reSensitive.ReplaceAllString(templateString, "{{ \"{{ ."+VAR+"."+m+" }}\" }}")
		}
	}

	tpl := createTemplate()
	tpl, err := tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: templateData.Constants,
		VAR:   templateData.Variables})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// ExecuteSensitiveTemplate templates the sensitive variables
// for use immediately before validation, after non-sensitive data is templated, results should not be written
func ExecuteSensitiveTemplate(templateData *TemplateData, templateString string) ([]byte, error) {
	tpl := createTemplate()
	tpl, err := tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	err = tpl.Execute(&buffer, map[string]interface{}{
		VAR: templateData.SensitiveVariables})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// ExecuteMaskedTemplate templates all values, but masks the sensitive ones
// for display/printing only
func ExecuteMaskedTemplate(templateData *TemplateData, templateString string) ([]byte, error) {
	// Find any sensitive keys {{ var.KEY }}, where KEY is in templateData.SensitiveVariables and replace with {{ var.KEY | mask }}
	re := regexp.MustCompile(`{{\s*\.` + VAR + `\.([a-zA-Z0-9_]+)\s*}}`)
	varMatches := re.FindStringSubmatch(templateString)
	for _, m := range varMatches {
		if _, ok := templateData.SensitiveVariables[m]; ok {
			reSensitive := regexp.MustCompile(`{{\s*\.` + VAR + `\.` + m + `\s*}}`)
			templateString = reSensitive.ReplaceAllString(templateString, "{{ ."+VAR+"."+m+" | mask }}")
		}
	}

	tpl := createTemplate()
	tpl, err := tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	allVars := MergeStringMaps(templateData.Variables, templateData.SensitiveVariables)
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: templateData.Constants,
		VAR:   allVars})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// Prepare the templateData object for use in templating
func CollectTemplatingData(constants map[string]interface{}, variables []VariableConfig) *TemplateData {
	// Create the TemplateData object from the constants and variables
	templateData := NewTemplateData()
	templateData.Constants = constants
	for _, variable := range variables {
		// convert '-' to '_' in the key and remove any special characters
		variable.Key = strings.ReplaceAll(variable.Key, "-", "_")
		re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
		variable.Key = re.ReplaceAllString(variable.Key, "")

		templateData.Variables[variable.Key] = variable.Default
		if variable.Sensitive {
			templateData.SensitiveVariables[variable.Key] = variable.Default
		}
	}

	// Get all environment variables with a specific prefix
	envMap := GetEnvVars(PREFIX)

	// Update the templateData with the environment variables overrides
	templateData.Variables = MergeStringMaps(templateData.Variables, envMap)
	templateData.SensitiveVariables = MergeStringMaps(templateData.SensitiveVariables, envMap)

	return templateData
}

// get all environment variables with the established prefix
func GetEnvVars(prefix string) map[string]string {
	envMap := make(map[string]string)

	// Get all environment variables
	envVars := os.Environ()

	// Iterate over environment variables
	for _, envVar := range envVars {
		// Split the environment variable into key and value
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		// Check if the key starts with the specified prefix
		if strings.HasPrefix(key, prefix) {
			// Remove the prefix from the key and convert to lowercase
			strippedKey := strings.TrimPrefix(key, prefix)
			envMap[strings.ToLower(strippedKey)] = value
		}
	}

	return envMap
}

// MergeStringMaps merges two maps of strings into a single map of strings.
// m2 will overwrite m1 if a key exists in both maps.
func MergeStringMaps(m1, m2 map[string]string) map[string]string {
	r := map[string]string{}

	for key, value := range m1 {
		r[key] = value
	}

	for key, value := range m2 {
		r[key] = value
	}

	return r
}

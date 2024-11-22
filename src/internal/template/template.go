package template

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

const (
	PREFIX = "LULA_VAR_"
	CONST  = "const"
	VAR    = "var"
)

type RenderType string

const (
	MASKED       RenderType = "masked"
	CONSTANTS    RenderType = "constants"
	NONSENSITIVE RenderType = "non-sensitive"
	ALL          RenderType = "all"
)

type TemplateRenderer struct {
	tpl          *template.Template
	templateData *TemplateData
}

func NewTemplateRenderer(templateData *TemplateData) *TemplateRenderer {
	return &TemplateRenderer{
		tpl:          createTemplate(),
		templateData: templateData,
	}
}

func (r *TemplateRenderer) Render(templateString string, t RenderType) ([]byte, error) {
	switch t {
	case MASKED:
		return r.ExecuteMaskedTemplate(templateString)
	case CONSTANTS:
		return r.ExecuteConstTemplate(templateString)
	case NONSENSITIVE:
		return r.ExecuteNonSensitiveTemplate(templateString)
	case ALL:
		return r.ExecuteFullTemplate(templateString)
	default:
		return []byte{}, fmt.Errorf("invalid render type: %s", t)
	}
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

type VariableConfig struct {
	Key       string
	Default   string
	Sensitive bool
}

// ExecuteFullTemplate templates everything
func (r *TemplateRenderer) ExecuteFullTemplate(templateString string) ([]byte, error) {
	tpl, err := r.tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	allVars := concatStringMaps(r.templateData.Variables, r.templateData.SensitiveVariables)
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: r.templateData.Constants,
		VAR:   allVars})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// ExecuteConstTemplate templates only constants
// this templates only values in the constants map
func (r *TemplateRenderer) ExecuteConstTemplate(templateString string) ([]byte, error) {
	// Find anything {{ var.KEY }} and replace with {{ "{{ var.KEY }}" }}
	re := regexp.MustCompile(`{{\s*\.` + VAR + `\.([a-zA-Z0-9_]+)\s*}}`)
	templateString = re.ReplaceAllString(templateString, "{{ \"{{ ."+VAR+".$1 }}\" }}")

	tpl, err := r.tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: r.templateData.Constants})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// ExecuteNonSensitiveTemplate templates only constants and non-sensitive variables
// used for compose operations
func (r *TemplateRenderer) ExecuteNonSensitiveTemplate(templateString string) ([]byte, error) {
	// Find any sensitive keys {{ var.KEY }}, where KEY is in templateData.SensitiveVariables and replace with {{ "{{ var.KEY }}" }}
	re := regexp.MustCompile(`{{\s*\.` + VAR + `\.([a-zA-Z0-9_]+)\s*}}`)
	varMatches := re.FindAllStringSubmatch(templateString, -1)
	uniqueMatches := returnUniqueMatches(varMatches, 1)
	for k, matches := range uniqueMatches {
		if _, ok := r.templateData.SensitiveVariables[matches[0]]; ok {
			templateString = strings.ReplaceAll(templateString, k, "{{ \""+k+"\" }}")
		}
	}

	tpl, err := r.tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: r.templateData.Constants,
		VAR:   r.templateData.Variables})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// ExecuteMaskedTemplate templates all values, but masks the sensitive ones
// for display/printing only
func (r *TemplateRenderer) ExecuteMaskedTemplate(templateString string) ([]byte, error) {
	// Find any sensitive keys {{ var.KEY }}, where KEY is in templateData.SensitiveVariables and replace with {{ var.KEY | mask }}
	re := regexp.MustCompile(`{{\s*\.` + VAR + `\.([a-zA-Z0-9_]+)\s*}}`)
	varMatches := re.FindAllStringSubmatch(templateString, -1)
	uniqueMatches := returnUniqueMatches(varMatches, 1)
	for k, matches := range uniqueMatches {
		if _, ok := r.templateData.SensitiveVariables[matches[0]]; ok {
			templateString = strings.ReplaceAll(templateString, k, "********")
		}
	}

	tpl, err := r.tpl.Parse(templateString)
	if err != nil {
		return []byte{}, err
	}

	var buffer strings.Builder
	err = tpl.Execute(&buffer, map[string]interface{}{
		CONST: r.templateData.Constants,
		VAR:   r.templateData.Variables})
	if err != nil {
		return []byte{}, err
	}

	return []byte(buffer.String()), nil
}

// Prepare the templateData object for use in templating
func CollectTemplatingData(constants map[string]interface{}, variables []VariableConfig, overrides map[string]string) (*TemplateData, error) {
	// Create the TemplateData object from the constants and variables
	templateData := NewTemplateData()

	// check for invalid characters in keys
	err := checkForInvalidKeys(constants, variables)
	if err != nil {
		return templateData, err
	}

	templateData.Constants = deepCopyMap(constants)

	for _, variable := range variables {
		if variable.Sensitive {
			templateData.SensitiveVariables[variable.Key] = variable.Default
		} else {
			templateData.Variables[variable.Key] = variable.Default
		}
	}

	// Get all environment variables with a specific prefix
	envMap := GetEnvVars(PREFIX)

	// Update the templateData with the environment variables overrides
	templateData.Variables = mergeStringMaps(templateData.Variables, envMap)
	templateData.SensitiveVariables = mergeStringMaps(templateData.SensitiveVariables, envMap)

	// Apply overrides
	err = overrideTemplateValues(templateData, overrides)
	if err != nil {
		return templateData, err
	}

	// Validate that all env vars have values - currently debug prints missing env vars (do we want to return an error?)
	var variablesMissing strings.Builder
	for k, v := range templateData.Variables {
		if v == "" {
			variablesMissing.WriteString(fmt.Sprintf("variable %s is missing a value;\n", k))
		}
	}
	for k, v := range templateData.SensitiveVariables {
		if v == "" {
			variablesMissing.WriteString(fmt.Sprintf("sensitive variable %s is missing a value;\n", k))
		}
	}
	message.Debug(variablesMissing.String())

	return templateData, nil
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

// IsTemplate checks if the given string contains valid template syntax
func IsTemplate(data string) bool {
	// Check for basic template syntax markers
	if !strings.Contains(data, "{{") || !strings.Contains(data, "}}") {
		return false
	}

	// Attempt to parse the template
	tpl := createTemplate()
	_, err := tpl.Parse(data)
	return err == nil
}

func ParseRenderType(item string) (RenderType, error) {
	switch strings.ToLower(item) {
	case "masked":
		return MASKED, nil
	case "constants":
		return CONSTANTS, nil
	case "non-sensitive":
		return NONSENSITIVE, nil
	case "all":
		return ALL, nil
	}
	return "", fmt.Errorf("invalid render type: %s", item)
}

// createTemplate creates a new template object
func createTemplate() *template.Template {
	// Register custom template functions
	funcMap := template.FuncMap{
		"concatToRegoList": func(a []any) string {
			return concatToRegoList(a)
		},
		"uuid":      newUUID,
		"timestamp": timestamp,
		// Add more custom functions as needed
	}

	// Parse the template and apply the function map
	tpl := template.New("template").Funcs(funcMap)
	tpl.Option("missingkey=error")

	return tpl
}

// mergeStringMaps merges two maps of strings into a single map of strings.
// m2 will overwrite m1 if a key exists in both maps, similar to left-join operation
func mergeStringMaps(m1, m2 map[string]string) map[string]string {
	r := map[string]string{}

	for key, value := range m1 {
		r[key] = value
	}

	for key, value := range m2 {
		// only add the key if it does exist in r
		if _, ok := r[key]; ok {
			r[key] = value
		}
	}

	return r
}

// concatStringMaps concatenates two maps of strings into a single map of strings.
// m2 will overwrite m1 if a key exists in both maps
func concatStringMaps(m1, m2 map[string]string) map[string]string {
	r := make(map[string]string)

	for key, value := range m1 {
		r[key] = value
	}

	for key, value := range m2 {
		r[key] = value
	}
	return r
}

// returnUniqueMatches returns a slice of unique matches from a slice of strings
func returnUniqueMatches(matches [][]string, captures int) map[string][]string {
	uniqueMatches := make(map[string][]string)
	for _, match := range matches {
		fullMatch := match[0]
		if _, exists := uniqueMatches[fullMatch]; !exists {
			uniqueMatches[fullMatch] = match[captures:]
		}
	}
	return uniqueMatches
}

// checkForInvalidKeys checks for invalid characters in keys for go text/template
// cannot contain '-' or '.'
func checkForInvalidKeys(constants map[string]interface{}, variables []VariableConfig) error {
	var errs strings.Builder

	containsInvalidChars := func(key string) {
		if strings.Contains(key, "-") {
			errs.WriteString(fmt.Sprintf("invalid key %s - cannot contain '-';", key))
		}
		if strings.Contains(key, ".") {
			errs.WriteString(fmt.Sprintf("invalid key %s - cannot contain '.';", key))
		}
	}

	// check for invalid characters in keys, recursively through constants
	var validateKeys func(m map[string]interface{})
	validateKeys = func(m map[string]interface{}) {
		for key, value := range m {
			containsInvalidChars(key)
			if nestedMap, ok := value.(map[string]interface{}); ok {
				validateKeys(nestedMap)
			}
		}
	}

	validateKeys(constants)

	// check for invalid characters in keys in variables
	for _, variable := range variables {
		containsInvalidChars(variable.Key)
	}

	if errs.Len() > 0 {
		return errors.New(errs.String()[:len(errs.String())-1])
	}

	return nil
}

// overrideTemplateValues overrides values in the templateData object with values from the overrides map
func overrideTemplateValues(templateData *TemplateData, overrides map[string]string) error {
	for path, value := range overrides {
		// for each key, check if .var or .const
		// if .var, set the value in the templateData.Variables or templateData.SensitiveVariables
		// if .const, set the value in the templateData.Constants
		if strings.HasPrefix(path, "."+VAR+".") {
			key := strings.TrimPrefix(path, "."+VAR+".")

			if _, ok := templateData.SensitiveVariables[key]; ok {
				templateData.SensitiveVariables[key] = value
			} else {
				templateData.Variables[key] = value
			}
		} else if strings.HasPrefix(path, "."+CONST+".") {
			// Set the value in the templateData.Constants
			key := strings.TrimPrefix(path, "."+CONST+".")
			err := setNestedValue(templateData.Constants, key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Helper function to set a value in a map based on a JSON-like key path
// Only supports basic map path (root.key.subkey)
func setNestedValue(m map[string]interface{}, path string, value interface{}) error {
	keys := strings.Split(path, ".")
	lastKey := keys[len(keys)-1]

	// Traverse the map, creating intermediate maps if necessary
	for _, key := range keys[:len(keys)-1] {
		if _, exists := m[key]; !exists {
			m[key] = make(map[string]interface{})
		}
		if nestedMap, ok := m[key].(map[string]interface{}); ok {
			m = nestedMap
		} else {
			return fmt.Errorf("path %s contains a non-map value", key)
		}
	}

	// Set the final value
	m[lastKey] = value
	return nil
}

func deepCopyMap(input map[string]interface{}) map[string]interface{} {
	if input == nil {
		return nil
	}

	// Create a new map to hold the copy
	copy := make(map[string]interface{})

	for key, value := range input {
		// Check the type of the value and copy accordingly
		switch v := value.(type) {
		case map[string]interface{}:
			// If the value is a map, recursively deep copy it
			copy[key] = deepCopyMap(v)
		case []interface{}:
			// If the value is a slice, deep copy each element
			copy[key] = deepCopySlice(v)
		default:
			// For other types (e.g., strings, ints), just assign directly
			copy[key] = v
		}
	}

	return copy
}

// Helper function to deep copy a slice of interface{}
func deepCopySlice(input []interface{}) []interface{} {
	if input == nil {
		return nil
	}

	// Create a new slice to hold the copy
	copy := make([]interface{}, len(input))

	for i, value := range input {
		// Check the type of the value and copy accordingly
		switch v := value.(type) {
		case map[string]interface{}:
			// If the value is a map, recursively deep copy it
			copy[i] = deepCopyMap(v)
		case []interface{}:
			// If the value is a slice, deep copy each element
			copy[i] = deepCopySlice(v)
		default:
			// For other types (e.g., strings, ints), just assign directly
			copy[i] = v
		}
	}

	return copy
}

package template_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/internal/template"
)

func testRender(t *testing.T, templateRenderer *template.TemplateRenderer, templateString string, renderType template.RenderType, expected string) error {
	t.Helper()

	got, err := templateRenderer.Render(templateString, renderType)
	if err != nil {
		return fmt.Errorf("error templating data: %v\n", err.Error())
	}

	if string(got) != expected {
		t.Fatalf("Expected %s - Got %s\n", expected, string(got))
	}
	return nil
}

func TestExecuteFullTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Test template all with data", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`
		expected := `
		constant template: testing
		variable template: my-env-var
		secret template: my-secret
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.ALL, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	// Note - this will change depending on the tpl.Option chosen
	t.Run("Test template all with all empty data, error", func(t *testing.T) {
		templateData := template.NewTemplateData()
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.ALL, "")
		if err == nil {
			t.Fatalf("Expected an error, but got nil")
		}
	})

	// Note - this will change depending on the tpl.Option chosen
	t.Run("Test template all with invalid template paths, error", func(t *testing.T) {
		templateData := template.NewTemplateData()

		templateString := `
		constant template: {{ .constant.testVar }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.ALL, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Test template all with invalid template characters, error", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"test-var": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}

		templateString := `
		constant template: {{ .const.test-var }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.ALL, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	// Note - this will change depending on the tpl.Option chosen
	t.Run("Test template all with invalid template subpath, error", func(t *testing.T) {
		templateData := template.NewTemplateData()

		templateString := `
		variable template: {{ .var.nokey.sub }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.ALL, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Test template all with invalid template, error", func(t *testing.T) {
		templateData := template.NewTemplateData()

		templateString := `
		constant template: {{ constant.testVar }}
		`
		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.ALL, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Test template on 'concatToRegoList' function", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"exemptions": []interface{}{"one", "two", "three"},
			},
		}

		templateString := `
		constant template: {{ .const.exemptions | concatToRegoList }}
		`
		expected := `
		constant template: "one", "two", "three"
		`
		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.ALL, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

}

func TestExecuteConstTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Test template const with data", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`
		expected := `
		constant template: testing
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.CONSTANTS, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test template const with missing var data", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`
		expected := `
		constant template: testing
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`
		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.CONSTANTS, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test template const with weird spacing in template", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{.var.some_env_var}}
		secret template: {{  .var.some_lula_secret  }}
		`
		expected := `
		constant template: testing
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`
		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.CONSTANTS, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	// Note - this will change depending on the tpl.Option chosen
	t.Run("Test template const with empty data", func(t *testing.T) {
		templateData := &template.TemplateData{
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.CONSTANTS, "")
		if err == nil {
			t.Fatalf("Expected an error, but got nil")
		}
	})

}

func TestExecuteNonSensitiveTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Test template nonsensitive with data and duplicate matches", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		variable template2: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		secret template2: {{ .var.some_lula_secret }}
		`
		expected := `
		constant template: testing
		variable template: my-env-var
		variable template2: my-env-var
		secret template: {{ .var.some_lula_secret }}
		secret template2: {{ .var.some_lula_secret }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.NONSENSITIVE, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test template nonsensitive with weird spacing in template", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{.var.some_lula_secret   }}
		`
		expected := `
		constant template: testing
		variable template: my-env-var
		secret template: {{.var.some_lula_secret   }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.NONSENSITIVE, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	// Note - this will change depending on the tpl.Option chosen
	t.Run("Test template nonsensitive with empty var data", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.NONSENSITIVE, "")
		if err == nil {
			t.Fatalf("Expected an error, but got nil")
		}
	})

}

func TestExecuteMaskedTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Test masked template with sensitive data and duplicates", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		variable template2: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		secret template2: {{ .var.some_lula_secret }}
		`
		expected := `
		constant template: testing
		variable template: my-env-var
		variable template2: my-env-var
		secret template: ********
		secret template2: ********
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.MASKED, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test masked template with weird spacing in template", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{.var.some_lula_secret   }}
		`
		expected := `
		constant template: testing
		variable template: my-env-var
		secret template: ********
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.MASKED, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	// Note - this will change depending on the tpl.Option chosen
	t.Run("Test masked template with missing data", func(t *testing.T) {
		templateData := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		templateString := `
		constant template: {{ .const.testVar }}
		variable template: {{ .var.some_env_var }}
		secret template: {{ .var.some_lula_secret }}
		`

		tr := template.NewTemplateRenderer(templateData)
		err := testRender(t, tr, templateString, template.MASKED, "")
		if err == nil {
			t.Fatalf("Expected an error, but got nil")
		}
	})
}

func TestCollectTemplatingData(t *testing.T) {

	test := func(t *testing.T, constants map[string]interface{}, variables []template.VariableConfig, overrides map[string]string, expected *template.TemplateData) error {
		t.Helper()
		// templateData returned
		got, err := template.CollectTemplatingData(constants, variables, overrides)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Expected %v - Got %v\n", expected, got)
		}
		return nil
	}

	t.Run("Test collect templating data", func(t *testing.T) {
		var overrides map[string]string
		constants := map[string]interface{}{
			"testVar": "testing",
		}
		variables := []template.VariableConfig{
			{
				Key:       "some_env_var",
				Default:   "my-env-var",
				Sensitive: false,
			},
			{
				Key:       "some_lula_secret",
				Default:   "my-secret",
				Sensitive: true,
			},
		}
		expected := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "my-secret",
			},
		}
		err := test(t, constants, variables, overrides, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test collect templating data with env vars", func(t *testing.T) {
		os.Setenv("LULA_VAR_SOME_LULA_SECRET", "env-secret")
		os.Setenv("LULA_VAR_SOME_ENV_VAR", "env-var")
		defer os.Unsetenv("LULA_VAR_SOME_LULA_SECRET")
		defer os.Unsetenv("LULA_VAR_SOME_ENV_VAR")

		var overrides map[string]string
		constants := map[string]interface{}{
			"testVar": "testing",
		}
		variables := []template.VariableConfig{
			{
				Key:       "some_env_var",
				Default:   "my-env-var",
				Sensitive: false,
			},
			{
				Key:       "some_lula_secret",
				Default:   "my-secret",
				Sensitive: true,
			},
		}
		expected := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "env-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "env-secret",
			},
		}
		err := test(t, constants, variables, overrides, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test collect templating data with overrides", func(t *testing.T) {
		overrides := map[string]string{
			".var.some_env_var":     "override-var",
			".var.some_lula_secret": "override-secret",
			".const.test.subkey":    "override-subkey",
		}
		constants := map[string]interface{}{
			"testVar": "testing",
			"test": map[string]interface{}{
				"subkey": "subkey-value",
			},
		}
		variables := []template.VariableConfig{
			{
				Key:       "some_env_var",
				Default:   "my-env-var",
				Sensitive: false,
			},
			{
				Key:       "some_lula_secret",
				Default:   "my-secret",
				Sensitive: true,
			},
		}
		expected := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
				"test": map[string]interface{}{
					"subkey": "override-subkey",
				},
			},
			Variables: map[string]string{
				"some_env_var": "override-var",
			},
			SensitiveVariables: map[string]string{
				"some_lula_secret": "override-secret",
			},
		}
		err := test(t, constants, variables, overrides, expected)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
	})

	t.Run("Test collect templating data with bad keys", func(t *testing.T) {
		var overrides map[string]string
		constants := map[string]interface{}{
			"test-var": "testing",
			"anotherkey": map[string]interface{}{
				"sub.key": "testing",
			},
		}
		variables := []template.VariableConfig{
			{
				Key:       "some-env-var",
				Default:   "my-env-var",
				Sensitive: false,
			},
		}
		expected := &template.TemplateData{
			Constants: map[string]interface{}{
				"testVar": "testing",
			},
			Variables: map[string]string{
				"some_env_var": "my-env-var",
			},
			SensitiveVariables: map[string]string{},
		}
		err := test(t, constants, variables, overrides, expected)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		numErrors := strings.Count(err.Error(), "invalid key")
		if numErrors != 3 {
			t.Fatalf("expected 3 invalid constant keys, got %d", numErrors)
		}
	})
}

func TestGetEnvVars(t *testing.T) {

	test := func(t *testing.T, prefix string, key string, value string) {
		t.Helper()

		os.Setenv(key, value)
		defer os.Unsetenv(key)
		envMap := template.GetEnvVars(prefix)

		// convert key to expected format
		strippedKey := strings.TrimPrefix(key, prefix)

		if envMap[strings.ToLower(strippedKey)] != value {
			t.Fatalf("Expected %s - Got %s\n", value, envMap[strings.ToLower(strippedKey)])
		}
	}

	t.Run("Test LULA_RESOURCE - Pass", func(t *testing.T) {
		test(t, "LULA_", "LULA_RESOURCE", "pods")
	})

	t.Run("Test OTHER_RESOURCE - Pass", func(t *testing.T) {
		test(t, "OTHER_", "OTHER_RESOURCE", "deployments")
	})
}

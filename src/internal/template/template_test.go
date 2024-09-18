package template_test

import (
	"os"
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/internal/template"
)

func TestExecuteTemplate(t *testing.T) {

	test := func(t *testing.T, data map[string]interface{}, preTemplate string, expected string) {
		t.Helper()
		// templateData returned
		got, err := template.ExecuteTemplate(data, preTemplate)
		if err != nil {
			t.Fatalf("error templating data: %s\n", err.Error())
		}

		if string(got) != expected {
			t.Fatalf("Expected %s - Got %s\n", expected, string(got))
		}
	}

	t.Run("Test {{ .testVar }} with data", func(t *testing.T) {
		data := map[string]interface{}{
			"testVar": "testing",
		}

		test(t, data, "{{ .testVar }}", "testing")
	})

	t.Run("Test {{ .testVar }} but empty data", func(t *testing.T) {
		data := map[string]interface{}{}

		test(t, data, "{{ .testVar }}", "<no value>")
	})

}

func TestGetEnvVars(t *testing.T) {

	test := func(t *testing.T, prefix string, key string, value string) {
		t.Helper()

		os.Setenv(key, value)
		envMap := template.GetEnvVars(prefix)

		// convert key to expected format
		strippedKey := strings.TrimPrefix(key, prefix)

		if envMap[strings.ToLower(strippedKey)] != value {
			t.Fatalf("Expected %s - Got %s\n", value, envMap[strings.ToLower(strippedKey)])
		}
		os.Unsetenv(key)
	}

	t.Run("Test LULA_RESOURCE - Pass", func(t *testing.T) {
		test(t, "LULA_", "LULA_RESOURCE", "pods")
	})

	t.Run("Test OTHER_RESOURCE - Pass", func(t *testing.T) {
		test(t, "OTHER_", "OTHER_RESOURCE", "deployments")
	})
}

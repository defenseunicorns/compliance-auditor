package template_test

import (
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

// func TestGetEnvVars(t *testing.T) {
// 	test := func(t *testing.T, data string, expected string) {
// 		t.Helper()
// 	}

// 	t.Run("Test One - Passing", func(t *testing.T) {
// 		test(t, "test", "test")
// 	})
// }

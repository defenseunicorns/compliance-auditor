package cmd_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/cmd"
	"github.com/defenseunicorns/lula/src/test/util"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestMain(m *testing.M) {
	flag.Parse()
	m.Run()
}

func runCmdTest(t *testing.T, goldenFileName string, expectError bool, cmdArgs ...string) error {
	t.Helper()

	rootCmd := cmd.RootCommand()

	_, output, err := util.ExecuteCommand(rootCmd, cmdArgs...)
	if err != nil {
		if !expectError {
			return err
		} else {
			return nil
		}
	}

	if !expectError {
		testGolden(t, goldenFileName, output)
	}

	return nil
}

func testGolden(t *testing.T, filename, got string) {
	t.Helper()

	got = strings.ReplaceAll(got, "\r\n", "\n")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	goldenPath := filepath.Join(wd, "testdata", filename+".golden")

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	wantBytes, _ := os.ReadFile(goldenPath)

	want := string(wantBytes)

	if got != want {
		t.Fatalf("`%s` does not match.\n\nWant:\n\n%s\n\nGot:\n\n%s", goldenPath, want, got)
	}
}

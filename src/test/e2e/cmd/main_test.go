package cmd_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func runCmdTestWithOutputFile(t *testing.T, goldenFileName string, outExt string, expectError bool, cmdArgs ...string) error {
	t.Helper()

	tempFileName := fmt.Sprintf("output-%s.%s", goldenFileName, outExt)
	defer os.Remove(tempFileName)

	rootCmd := cmd.RootCommand()

	cmdArgs = append(cmdArgs, "-o", tempFileName)
	_, _, err := util.ExecuteCommand(rootCmd, cmdArgs...)
	if err != nil {
		if !expectError {
			return err
		} else {
			return nil
		}
	}

	// Read the output file
	data, err := os.ReadFile(tempFileName)
	if err != nil {
		return err
	}

	// Scrub timestamps
	data = scrubTimestamps(data)

	if !expectError {
		testGolden(t, goldenFileName, string(data))
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

func scrubTimestamps(data []byte) []byte {
	re := regexp.MustCompile(`(?i)(last-modified:\s*)(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+[-+]\d{2}:\d{2})`)
	return []byte(re.ReplaceAllString(string(data), "${1}XXX"))
}

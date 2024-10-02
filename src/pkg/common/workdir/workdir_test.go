package workdir_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/defenseunicorns/lula/src/pkg/common/workdir"
)

func TestSwitchCwd(t *testing.T) {

	tempDir := t.TempDir()

	tests := []struct {
		name     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid path",
			path:     tempDir,
			expected: tempDir,
			wantErr:  false,
		},
		{
			name:     "Path is File",
			path:     "./workdir_test.go",
			expected: "./",
			wantErr:  false,
		},
		{
			name:     "Invalid path",
			path:     "/nonexistent",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty Path",
			path:     "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFunc, err := workdir.SetCwdToFileDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SwitchCwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				defer resetFunc()
				wd, _ := os.Getwd()
				expected, err := filepath.Abs(tt.expected)
				if err != nil {
					t.Errorf("SwitchCwd() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !strings.HasSuffix(wd, expected) {
					t.Errorf("SwitchCwd() working directory = %v, want %v", wd, tt.expected)
				}
			}
		})
	}
}

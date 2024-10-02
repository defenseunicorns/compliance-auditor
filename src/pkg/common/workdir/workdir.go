package workdir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/lula/src/pkg/message"
)

// SetCwdToFileDir takes a path and changes the current working directory to the directory of the path
func SetCwdToFileDir(dirPath string) (resetFunc func(), err error) {
	// Bail if empty
	if dirPath == "" {
		return resetFunc, fmt.Errorf("dirPath is empty")
	}
	info, err := os.Stat(dirPath)
	// Bail if the path does not exist
	if err != nil {
		return resetFunc, err
	}
	// Get the directory of the path if it is not a directory
	if !info.IsDir() {
		dirPath = filepath.Dir(dirPath)
	}

	// Save the current working directory
	originalDir, err := os.Getwd()
	if err != nil {
		return resetFunc, err
	}

	// Change the current working directory
	err = os.Chdir(dirPath)
	if err != nil {
		return resetFunc, err
	}

	// Return a function to change the current working directory back to the original directory
	resetFunc = func() {
		err = os.Chdir(originalDir)
		if err != nil {
			message.Warnf("unable to change cwd back to %s: %v", originalDir, err)
		}
	}

	// Change back to the original working directory when done
	return resetFunc, err
}

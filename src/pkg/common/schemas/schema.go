package schemas

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed *.json
var Schemas embed.FS

const (
	SCHEMA_SUFFIX = ".json"
)

// HasSchema checks if a schema exists in the schemas directory
func HasSchema(path string) bool {
	if !strings.HasSuffix(path, SCHEMA_SUFFIX) {
		path = path + SCHEMA_SUFFIX
	}
	_, err := Schemas.Open(path)
	return err == nil
}

func ListSchemas() ([]string, error) {
	files, err := ToMap()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	return keys, nil
}

// ToMap returns a map of schema names to schemas
func ToMap() (fileMap map[string]fs.DirEntry, err error) {
	files, err := Schemas.ReadDir(".")
	if err != nil {
		return nil, err
	}
	fileMap = make(map[string]fs.DirEntry)
	for _, file := range files {
		name := file.Name()
		isDir := file.IsDir()
		if isDir || !strings.HasSuffix(name, SCHEMA_SUFFIX) {
			continue
		}
		fileMap[name] = file
	}
	return fileMap, nil
}

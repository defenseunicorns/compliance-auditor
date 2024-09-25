package common

import (
	"fmt"
	"strings"

	"github.com/defenseunicorns/lula/src/internal/template"
)

func ParseTemplateOverrides(setFlags []string) (map[string]string, error) {
	overrides := make(map[string]string)
	for _, flag := range setFlags {
		parts := strings.SplitN(flag, "=", 2)
		if len(parts) != 2 {
			return overrides, fmt.Errorf("invalid --set flag format, should be .root.key=value")
		}

		if !strings.HasPrefix(parts[0], "."+template.CONST+".") && !strings.HasPrefix(parts[0], "."+template.VAR+".") {
			return overrides, fmt.Errorf("invalid --set flag format, path should start with .const or .var")
		}

		path, value := parts[0], parts[1]
		overrides[path] = value
	}
	return overrides, nil
}

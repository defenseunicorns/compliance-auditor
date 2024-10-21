package api

import (
	"errors"
	"fmt"
)

func validateSpec(spec *ApiSpec) error {
	if spec == nil {
		return errors.New("spec is required")
	}
	if len(spec.Requests) == 0 {
		return errors.New("some requests must be specified")
	}
	for _, request := range spec.Requests {
		if request.Name == "" {
			return errors.New("request name cannot be empty")
		}
		if request.URL == "" {
			return errors.New("request url cannot be empty")
		}
		if request.Method != "" {
			if request.Method != "Get" && request.Method != "Head" && request.Method != "Post" && request.Method != "PostForm" {
				return fmt.Errorf("unsupported method: %s", request.Method)
			}
		}
	}
	return nil
}

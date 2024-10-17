package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/defenseunicorns/lula/src/types"
)

func MakeRequests(Requests []Request) (types.DomainResources, error) {
	collection := make(map[string]interface{}, 0)

	for _, request := range Requests {
		transport := &http.Transport{}
		client := &http.Client{Transport: transport}

		resp, err := client.Get(request.URL)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil,
				fmt.Errorf("expected status code 200 but got %d", resp.StatusCode)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType == "application/json" {

			var prettyBuff bytes.Buffer
			err := json.Indent(&prettyBuff, body, "", "  ")
			if err != nil {
				return nil, err
			}
			prettyJson := prettyBuff.String()

			var tempData interface{}
			err = json.Unmarshal([]byte(prettyJson), &tempData)
			if err != nil {
				return nil, err
			}
			collection[request.Name] = tempData

		} else {
			return nil, fmt.Errorf("content type %s is not supported", contentType)
		}
	}
	return collection, nil
}

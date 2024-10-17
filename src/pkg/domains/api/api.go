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
<<<<<<< HEAD
			_ = json.Indent(&prettyBuff, body, "", "  ") // G104
=======
			err := json.Indent(&prettyBuff, body, "", "  ")
			if err != nil {
				return nil, err
			}
>>>>>>> 05a7f6e075adafeb2b3c635803cc12083f12c01d
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

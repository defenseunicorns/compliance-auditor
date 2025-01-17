package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"

	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/types"
)

type APIResponse struct {
	StatusCode int
	Status     string
	Raw        json.RawMessage
	Response   any
}

func (a ApiDomain) makeRequests(ctx context.Context) (types.DomainResources, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("canceled: %s", ctx.Err())
	default:
		collection := make(map[string]interface{}, 0)

		// defaultOpts apply to all requests, but may be overridden by adding an
		// options block to an individual request.
		var defaultOpts *ApiOpts
		if a.Spec.Options == nil {
			// This isn't likely to be nil in real usage, since CreateApiDomain
			// parses and mutates specs.
			defaultOpts = new(ApiOpts)
			defaultOpts.timeout = &defaultTimeout
		} else {
			defaultOpts = a.Spec.Options
		}

		// configure the default HTTP client using any top-level Options. Individual
		// requests with overrides (in request.Options.Headers) will get bespoke clients.
		defaultClient := clientFromOpts(defaultOpts)
		var errs error
		for i, request := range a.Spec.Requests {
			if i > 0 { // the first request cannot use outputs
				if a.Spec.outputs != nil {
					var err error
					request, err = executeTpls(request, a.Spec.outputs)
					if err != nil {
						/// TODO: better error
						return collection, err
					}
				}
			}

			var r io.Reader
			if request.Body != "" {
				r = bytes.NewBufferString(request.Body)
			}

			var headers map[string]string
			var client http.Client

			if request.Options == nil {
				headers = defaultOpts.Headers
				client = defaultClient
			} else {
				headers = request.Options.Headers
				client = clientFromOpts(request.Options)
			}

			response, err := doHTTPReq(ctx, client, request.Method, *request.reqURL, r, headers, request.reqParameters)
			if err != nil {
				errs = errors.Join(errs, err)
			}
			if response != nil {
				collection[request.Name] = types.DomainResources{
					"status":     response.Status,
					"statuscode": response.StatusCode,
					"raw":        response.Raw,
					"response":   response.Response,
				}

				if request.Outputs != nil {
					if a.Spec.outputs == nil {
						a.Spec.outputs = make(map[string]map[string]interface{})
					}
					message.Debug("processing outputs")
					node, err := yaml.ConvertJSONToYamlNode(string(response.Raw))
					if err != nil {
						errs = errors.Join(errs, err)
						break
					}
					for _, output := range request.Outputs {
						v, err := node.GetFieldValue(output.Path)
						if err != nil {
							errs = errors.Join(errs, err)
						}
						a.Spec.outputs[request.Name] = make(map[string]interface{})
						a.Spec.outputs[request.Name][output.Name] = v
					}
				}
			} else {
				// If the entire response is empty, return a validly empty resource
				collection[request.Name] = types.DomainResources{"status": 0}
			}
		}
		return collection, errs
	}
}

func executeTpls(req Request, vars map[string]map[string]interface{}) (Request, error) {
	modifiedReq := req.DeepCopy()

	// url
	if req.URLTpl != "" {
		urlstr, err := executeTpl(req.URLTpl, vars)
		if err != nil {
			return modifiedReq, err
		}
		reqUrl, err := url.Parse(urlstr)
		if err != nil {
			// TODO: error types! this is a copy!
			// TODO: we should attempt each template step and return a wrapped error
			return modifiedReq, errors.New("invalid request url")
		}
		modifiedReq.reqURL = reqUrl
	}

	// headers
	if req.Options.HeadersTpl != nil {
		message.Debug("overriding headers")
		headers := make(map[string]string)
		for k, v := range req.Options.HeadersTpl {
			message.Debug("executing tpl %v with vars %v", v, vars)
			h, err := executeTpl(v, vars)
			if err != nil {
				return modifiedReq, err
			}
			headers[k] = h
		}
		modifiedReq.Options.Headers = headers
	}

	// params
	if req.ParamsTpl != nil {
		params := make(map[string]string)
		for k, v := range req.ParamsTpl {
			h, err := executeTpl(v, vars)
			if err != nil {
				return modifiedReq, err
			}
			params[k] = h
		}
		modifiedReq.Params = params // this isn't strictly _necessary_ but let's get all the info in place
		queryParameters := url.Values{}
		for k, v := range modifiedReq.Params {
			queryParameters.Add(k, v)
		}
		modifiedReq.reqParameters = queryParameters
	}

	return modifiedReq, nil
}

func executeTpl(tplstr string, vars map[string]map[string]interface{}) (string, error) {
	tpl, err := template.New("lula").Delims("[[", "]]").Parse(tplstr)
	if err != nil {
		return "", err
	}

	var res bytes.Buffer
	err = tpl.Execute(&res, vars)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

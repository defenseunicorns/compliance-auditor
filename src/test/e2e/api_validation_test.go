package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/defenseunicorns/lula/src/internal/template"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/validation"
	"github.com/defenseunicorns/lula/src/pkg/message"
	"github.com/defenseunicorns/lula/src/test/util"
)

func TestApiValidation(t *testing.T) {
	const (
		ckAPIFieldConfigMap contextKey = "api-field-configmap"
		ckApiFieldPod       contextKey = "api-field-pod"
	)
	featureTrueValidation := features.New("Check API Validation - Success").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the configmap
			configMap, err := util.GetConfigMap("./scenarios/api-field/configmap.pass.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckAPIFieldConfigMap, configMap)

			// Create the pod
			pod, err := util.GetPod("./scenarios/api-field/pod.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckApiFieldPod, pod)

			return ctx
		}).
		Assess("Validate API response field", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/api-field/oscal-component.yaml"
			message.NoProgress = true

			validator, err := validation.New()
			if err != nil {
				t.Errorf("error creating validation context: %v", err)
			}

			assessment, err := validator.ValidateOnPath(context.Background(), oscalPath, "")
			if err != nil {
				t.Fatal(err)
			}

			if len(assessment.Results) == 0 {
				t.Fatal("Expected greater than zero results")
			}

			result := assessment.Results[0]

			if result.Findings == nil {
				t.Fatal("Expected findings to be not nil")
			}

			for _, finding := range *result.Findings {
				state := finding.Target.Status.State
				if state != "satisfied" {
					t.Fatal("State should be satisfied, but got :", state)
				}
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value(ckApiFieldPod).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err := wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(pod),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			configMap := ctx.Value(ckAPIFieldConfigMap).(*corev1.ConfigMap)
			if err := config.Client().Resources().Delete(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(configMap),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()

	featureFalseValidation := features.New("Check API Validation - Failure").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			configMap, err := util.GetConfigMap("./scenarios/api-field/configmap.fail.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckAPIFieldConfigMap, configMap)

			pod, err := util.GetPod("./scenarios/api-field/pod.yaml")
			if err != nil {
				t.Fatal(err)
			}
			if err = config.Client().Resources().Create(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					PodConditionMatch(pod, corev1.PodReady, corev1.ConditionTrue),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}
			ctx = context.WithValue(ctx, ckApiFieldPod, pod)
			return ctx
		}).
		Assess("Validate API response field", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			oscalPath := "./scenarios/api-field/oscal-component.yaml"
			message.NoProgress = true

			validator, err := validation.New()
			if err != nil {
				t.Errorf("error creating validation context: %v", err)
			}

			assessment, err := validator.ValidateOnPath(context.Background(), oscalPath, "")
			if err != nil {
				t.Fatal(err)
			}

			if len(assessment.Results) == 0 {
				t.Fatal("Expected greater than zero results")
			}

			result := assessment.Results[0]

			if result.Findings == nil {
				t.Fatal("Expected findings to be not nil")
			}

			for _, finding := range *result.Findings {
				state := finding.Target.Status.State
				if state != "not-satisfied" {
					t.Fatal("State should be not-satisfied, but got :", state)
				}
			}

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			pod := ctx.Value(ckApiFieldPod).(*corev1.Pod)
			if err := config.Client().Resources().Delete(ctx, pod); err != nil {
				t.Fatal(err)
			}
			err := wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(pod),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			configMap := ctx.Value(ckAPIFieldConfigMap).(*corev1.ConfigMap)
			if err := config.Client().Resources().Delete(ctx, configMap); err != nil {
				t.Fatal(err)
			}
			err = wait.
				For(conditions.New(config.Client().Resources()).
					ResourceDeleted(configMap),
					wait.WithTimeout(time.Minute*5))
			if err != nil {
				t.Fatal(err)
			}

			return ctx
		}).Feature()

	testEnv.Test(t, featureTrueValidation, featureFalseValidation)
}

// TestApiValidation_templated uses a URL parameter to control the return response from the API.
func TestApiValidation_templatedGet(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for the custom header included in the Request
		gotHeader := r.Header.Get("x-special-header")
		require.Equal(t, gotHeader, "lula")

		w.Header().Set("Content-Type", "application/json")
		wantResp := r.URL.Query().Get("response")
		require.NotEmpty(t, wantResp)
		passRsp := false
		if wantResp == "true" {
			passRsp = true
		}
		resp := struct {
			Pass bool `json:"pass"`
		}{
			passRsp,
		}
		err := json.NewEncoder(w).Encode(resp)
		require.NoError(t, err)
	}))
	defer svr.Close()

	tmpl := "scenarios/api-validations/component-definition.yaml.tmpl"

	// since it's just the two tests I'm using the name to check the assessment result.
	tests := map[string]struct {
		response string
	}{
		"satisfied":     {"true"},
		"not-satisfied": {"false"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			composer, err := composition.New(
				composition.WithModelFromLocalPath(tmpl),
				composition.WithRenderSettings("all", true),
				composition.WithTemplateRenderer("all", nil, []template.VariableConfig{
					{
						Key:     "reqUrl",
						Default: svr.URL,
					},
					{
						Key:     "response",
						Default: test.response,
					},
				}, []string{}),
			)
			require.NoError(t, err)

			validator, err := validation.New(validation.WithComposition(composer, tmpl))
			require.NoError(t, err)

			assessment, err := validator.ValidateOnPath(context.Background(), tmpl, "")
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(assessment.Results), 1)

			result := assessment.Results[0]
			require.NotNil(t, result.Findings)
			for _, finding := range *result.Findings {
				state := finding.Target.Status.State
				if state != name {
					t.Fatalf("State should be %s, but got %s", name, state)
				}
			}
		})
	}
}

// Similar to above, but using Post + Body instead of Get + Parameters
func TestApiValidation_templatedPost(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//there's a header in the request
		gotHeader := r.Header.Get("x-special-header")
		require.Equal(t, gotHeader, "lula")

		w.Header().Set("Content-Type", "application/json")
		wantResp, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.NotEmpty(t, wantResp)

		passRsp := false
		if string(wantResp) == "true\n" { //am I reading this badly? Or sending it badly? either way this is annoying.
			passRsp = true
		}
		resp := struct {
			Pass bool `json:"pass"`
		}{
			passRsp,
		}

		err = json.NewEncoder(w).Encode(resp)
		require.NoError(t, err)
	}))
	defer svr.Close()

	tmpl := "scenarios/api-validations/component-definition-post.yaml.tmpl"

	// since it's just the two tests I'm using the name to check the assessment result.
	tests := map[string]struct {
		body string
	}{
		"satisfied":     {"true"},
		"not-satisfied": {"false"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			composer, err := composition.New(
				composition.WithModelFromLocalPath(tmpl),
				composition.WithRenderSettings("all", true),
				composition.WithTemplateRenderer("all", nil, []template.VariableConfig{
					{
						Key:     "reqUrl",
						Default: svr.URL,
					},
					{
						Key:     "body",
						Default: test.body,
					},
				}, []string{}),
			)
			require.NoError(t, err)

			// Let's prove that the API DOmain is picking up the `executable`
			// flag from the request in the test fixture by forcing
			// AllowExecution to false.
			validator, err := validation.New(validation.WithComposition(composer, tmpl), validation.WithAllowExecution(false, true))
			require.NoError(t, err)

			_, err = validator.ValidateOnPath(context.Background(), tmpl, "")
			// we don't return an error when the validation doesn't run, so the best we can do is check that there was no error.
			require.NoError(t, err)

			// if we refactor to return the validation error when execution is not allowed, we can check for this error - but a number of other tests fail
			//require.ErrorIs(t, err, types.ErrExecutionNotAllowed) // Validations requiring execution will NOT be run

			// ok, now that's failed as expected:
			validator, err = validation.New(validation.WithComposition(composer, tmpl), validation.WithAllowExecution(true, true))
			require.NoError(t, err)

			assessment, err := validator.ValidateOnPath(context.Background(), tmpl, "")
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(assessment.Results), 1)

			result := assessment.Results[0]
			require.NotNil(t, result.Findings)
			for _, finding := range *result.Findings {
				state := finding.Target.Status.State
				if state != name {
					t.Fatalf("State should be %s, but got %s", name, state)
				}
			}
		})
	}
}

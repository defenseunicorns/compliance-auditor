package api

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteTpl(t *testing.T) {
	tests := map[string]struct {
		tplStr    string
		vars      map[string]map[string]interface{}
		want      string
		expectErr bool
	}{
		"basic success": {
			tplStr: "hello, [[.login.username]]!",
			vars: map[string]map[string]interface{}{
				"login": {
					"username": "cheetarah",
				},
			},
			want:      "hello, cheetarah!",
			expectErr: false,
		},
		"success": {
			tplStr: "hello, [[ (index .login.usernames 0) ]]!",
			vars: map[string]map[string]interface{}{
				"login": {
					"usernames": []string{"batcat", "cheetarah"},
				},
			},
			want:      "hello, batcat!",
			expectErr: false,
		},
		//TODO: we should catch this at an earlier stage (validation schema)
		"invalid tplStr, not an error": {
			tplStr: "hello, {{.login.username]]!",
			vars: map[string]map[string]interface{}{
				"login": {
					"username": "cheetarah",
				},
			},
			want:      "",
			expectErr: false,
		},
		// i dunno what makes for an error, since the above example didn't do it
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := executeTpl(test.tplStr, test.vars)
			if (err != nil) != test.expectErr {
				t.Fatalf("error = %v, expectErr %v", err, test.expectErr)
			}
			if err != nil {
				require.Equal(t, test.want, got)
			}
		})
	}
}

func TestExecuteTpls(t *testing.T) {
	input := Request{
		ParamsTpl: map[string]string{
			"param1": "[[.input.value1]]",
			"param2": "[[.input.value2]]",
		},
		URLTpl: `http:example.com/[[.admin.realm]]/foo`,
		Options: &ApiOpts{
			HeadersTpl: map[string]string{
				"token": "[[.auth.token]]",
			},
		},
	}
	vars := map[string]map[string]interface{}{
		"auth": {
			"token": "s00p3rs3kr3t",
		},
		"input": {
			"value1": "foo",
			"value2": 3,
		},
		"admin": {
			"realm": "uds",
		},
	}
	request, err := executeTpls(input, vars)
	require.NoError(t, err)
	wantUrl, err := url.Parse("http:example.com/uds/foo")
	require.NoError(t, err)
	require.Equal(t, wantUrl, request.reqURL)
	require.Equal(t, map[string]string{
		"param1": "foo",
		"param2": "3",
	}, request.Params)
	require.Equal(t, map[string]string{
		"token": "s00p3rs3kr3t",
	}, request.Options.Headers)
}

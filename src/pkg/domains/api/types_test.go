package api_test

import (
	"testing"

	api "github.com/defenseunicorns/lula/src/pkg/domains/api"
)

func TestCreateApiDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		spec        *api.ApiSpec
		expectedErr bool
	}{
		{
			name:        "nil spec",
			spec:        nil,
			expectedErr: true,
		},
		{
			name: "empty requests",
			spec: &api.ApiSpec{
				Requests: []api.Request{},
			},
			expectedErr: true,
		},
		{
			name: "invalid request - no name",
			spec: &api.ApiSpec{
				Requests: []api.Request{
					{
						URL: "test",
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid request - no url",
			spec: &api.ApiSpec{
				Requests: []api.Request{
					{
						Name: "test",
					},
				},
			},
			expectedErr: true,
		},
		{
			name: "valid request",
			spec: &api.ApiSpec{
				Requests: []api.Request{
					{
						Name: "test",
						URL:  "test",
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := api.CreateApiDomain(tt.spec)
			if (err != nil) != tt.expectedErr {
				t.Errorf("CreateApiDomain() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
		})
	}
}

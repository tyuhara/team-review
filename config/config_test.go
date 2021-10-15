package config

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	t.Parallel()

	cases := []struct {
		envvars map[string]string
		isErr   bool
		want    *Config
	}{
		{
			map[string]string{
				"LOG_LEVEL":              "INFO",
				"GITHUB_APP_ID":          "123456",
				"GITHUB_APP_PRIVATE_KEY": "private-key",
				"GITHUB_APP_SECRET":      "secret",
			},
			false,
			&Config{
				LogLevel:            zapcore.InfoLevel,
				GithubAppID:         123456,
				GithubAppPrivateKey: "private-key",
				GithubAppSecret:     "secret",
			},
		},
		{
			map[string]string{
				"LOG_LEVEL":     "INFO",
				"GITHUB_APP_ID": "123456",
			},
			true,
			nil,
		},
	}

	for _, tc := range cases {
		var got Config
		err := envconfig.ProcessWith(context.Background(), &got, envconfig.MapLookuper(tc.envvars))

		if tc.isErr {
			if err == nil {
				t.Errorf("this case should be error: %v", tc.envvars)
			}
			continue
		}

		if err != nil {
			t.Errorf("failed: %v", err)
		}

		if diff := cmp.Diff(&got, tc.want); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	}
}

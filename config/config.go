package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	LogLevel               zapcore.Level `env:"LOG_LEVEL,default=Info"`
	GithubAppID            int64         `env:"GITHUB_APP_ID,required"`
	GithubAppPrivateKey    string        `env:"GITHUB_APP_PRIVATE_KEY,required"`
	GithubAppSecret        string        `env:"GITHUB_APP_SECRET,required"`
	GithubBranchProtection string        `env:"GITHUB_BRANCH_PROTECTION"`
}

//New returns type of Config
func New(ctx context.Context) (*Config, error) {
	var c Config
	err := envconfig.Process(ctx, &c)
	if err != nil {
		return &Config{}, err
	}

	return &c, nil
}

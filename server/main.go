package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tyuhara/team-review/config"
	"github.com/tyuhara/team-review/http"
	"github.com/tyuhara/team-review/log"
)

func main() {
	conf, err := config.New(context.Background())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Failed to read environment variables: %s\n", err)
		return
	}

	logger, err := log.New(conf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[ERROR] Failed to setup logger: %s\n", err)
		return
	}
	defer func() {
		_ = logger.Sync()
	}()
	sugar := logger.Sugar()

	http.Handler(conf, sugar)
}

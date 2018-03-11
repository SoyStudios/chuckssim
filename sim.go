package main

import (
	"os"

	"github.com/go-kit/kit/log"
)

const appName = "chuckssim"

var AppVersion = "dev"
var AppBuildID = "HEAD"

var (
	logger log.Logger
	cfg    *config
)

func main() {
	logger = log.NewSyncLogger(log.NewLogfmtLogger(os.Stdout))
	logger = log.WithPrefix(logger,
		"t", log.DefaultTimestampUTC,
		"appName", appName,
		"appVersion", AppVersion,
		"appBuild", AppBuildID,
		"caller", log.DefaultCaller,
	)
	// nolint: errcheck
	logger.Log("level", "info",
		"msg", "starting up")
	var err error
	cfg, err = newDefaultConfig()
	if err != nil {
		// nolint: errcheck
		logger.Log("level", "error",
			"msg", "error initializing config",
			"err", err)
		os.Exit(1)
	}
}

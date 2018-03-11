package main

import (
	"os"

	"chuckssim.soystudios.com/chuckssim/pkg/simulation"
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

	sim, err := simulation.New()
	if err != nil {
		logger.Log("level", "error",
			"msg", "error initializing simulation",
			"err", err,
		)
		os.Exit(1)
	}
}

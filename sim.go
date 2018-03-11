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
	logger.Log("level", "info",
		"msg", "starting up")
}

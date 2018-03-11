package main

import (
	"context"
	"net/http"
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
	ctx    context.Context
	cancel context.CancelFunc
)

func main() {
	ctx = context.Background()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

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
		// nolint: errcheck
		logger.Log("level", "error",
			"msg", "error initializing simulation",
			"err", err,
		)
		os.Exit(1)
	}
	err = sim.GenerateRandom(100, 100, 20, 50)
	if err != nil {
		// nolint: errcheck
		logger.Log("level", "error",
			"msg", "error generating random simulation",
			"err", err,
		)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr: cfg.address,

		ReadTimeout:  cfg.readTimeout,
		WriteTimeout: cfg.writeTimeout,
		IdleTimeout:  cfg.idleTimeout,
	}
	err = srv.ListenAndServe()
	if err != nil {
		// nolint: errcheck
		logger.Log("level", "error",
			"msg", "error serviong",
			"err", err,
		)
		os.Exit(1)
	}
}

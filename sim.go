package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"chuckssim.soystudios.com/chuckssim/pkg/bot"
	"chuckssim.soystudios.com/chuckssim/pkg/simulation"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/websocket"
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

type (
	BotDetail struct {
		Bot bot.Bot `json:"bot"`

		Type string `json:"type"`
	}
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
	cfg.parseFlags()

	sim, err := simulation.New()
	if err != nil {
		// nolint: errcheck
		logger.Log("level", "error",
			"msg", "error initializing simulation",
			"err", err,
		)
		os.Exit(1)
	}
	sim.BotSize = 20
	err = sim.GenerateRandom(1000, 1000, 20, 50)
	if err != nil {
		// nolint: errcheck
		logger.Log("level", "error",
			"msg", "error generating random simulation",
			"err", err,
		)
		os.Exit(1)
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		EnableCompression: true,

		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}
	srv := &http.Server{
		Addr: cfg.address,

		Handler: handler(upgrader, sim, logger),

		ReadTimeout:  cfg.readTimeout,
		WriteTimeout: cfg.writeTimeout,
		IdleTimeout:  cfg.idleTimeout,
	}
	// nolint: errcheck
	logger.Log("level", "info",
		"msg", "serving HTTP",
		"address", srv.Addr,
	)
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

func handler(upgrader websocket.Upgrader, sim *simulation.Simulation, logger log.Logger) http.Handler {
	logger = log.With(logger,
		"transport", "HTTP",
		"context", "Websocket",
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// nolint: errcheck
			logger.Log("level", "error",
				"msg", "error upgrading WS",
				"err", err,
			)
			return
		}
		go serveSimulation(sim, conn, logger)
	})
}

func serveSimulation(sim *simulation.Simulation, conn *websocket.Conn, logger log.Logger) {
	ticker := time.NewTicker(1000 / 24 * time.Millisecond)
	defer func() {
		ticker.Stop()
		// nolint: errcheck
		conn.Close()
	}()
	messages := make(chan []byte)
	go func() {
		defer func() {
			close(messages)
		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseNoStatusReceived,
					websocket.CloseGoingAway,
				) {
					return
				}
				// nolint: errcheck
				logger.Log("level", "error",
					"msg", "error reading ws message",
					"err", err)
			}
			messages <- msg
		}
	}()
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				return
			}
			// nolint: errcheck
			logger.Log("level", "info",
				"msg", "received WS msg",
				"content", msg)
			if strings.HasPrefix(string(msg), "bot/") {
				logger.Log("level", "info", "msg", "bot info requested")
				msgList := strings.Split(string(msg), "/")
				id, err := strconv.ParseInt(msgList[1], 10, 64)
				if err != nil {
					return
				}
				bot, ok := sim.Bots[id]
				if !ok {
					continue
				}
				logger.Log("level", "info", "msg", bot.ID)
				b := BotDetail{Type: "BotDetail", Bot: bot}
				err = conn.WriteJSON(b)
				if err != nil {
					// nolint: errcheck
					logger.Log("level", "error",
						"msg", "error writing bot",
						"err", err)
					return
				}
			}
		case <-ticker.C:
			err := conn.WriteJSON(sim)
			if err != nil {
				// nolint: errcheck
				logger.Log("level", "error",
					"msg", "error writing state",
					"err", err)
				return
			}
		}
	}
}

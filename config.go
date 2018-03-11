package main

import "time"

type config struct {
	address string

	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

func newDefaultConfig() (*config, error) {
	cfg := &config{}

	cfg.address = "127.0.0.1:8080"

	cfg.readTimeout = 5 * time.Second
	cfg.writeTimeout = 10 * time.Second
	cfg.idleTimeout = 2 * time.Minute

	return cfg, nil
}

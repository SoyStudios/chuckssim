package main

type config struct {
	address string
}

func newDefaultConfig() (*config, error) {
	cfg := &config{}

	cfg.address = "127.0.0.1:8080"

	return cfg, nil
}

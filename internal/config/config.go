package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgresql://root@localhost:26257/mini_ledger?sslmode=disable"`
	HTTPPort    string `env:"HTTP_PORT" envDefault:"8080"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
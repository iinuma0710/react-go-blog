package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	BackendEnv string `env:"BACKEND_ENV" envDefault:"dev"`
	BckendPort int    `env:"BACKEND_PORT" envDefault:"80"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

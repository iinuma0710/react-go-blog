package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	BackendEnv string `env:"BACKEND_ENV" envDefault:"dev"`
	BckendPort int    `env:"BACKEND_PORT" envDefault:"80"`
	DBHost     string `env:"BLOG_DATABASE_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"BLOG_DATABASE_PORT" envDefault:"3306"`
	DBUser     string `env:"BLOG_DATABASE_USER" envDefault:"blog"`
	DBPassword string `env:"BLOG_DATABASE_PASSWORD" envDefault:"blog"`
	DBName     string `env:"BLOG_DATABASE_DATABASE" envDefault:"blog"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

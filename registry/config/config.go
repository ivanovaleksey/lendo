package config

import (
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/nats"
	"github.com/ivanovaleksey/lendo/registry/bank"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Bank bank.Config `envconfig:"bank"`
	DB   db.Config   `envconfig:"db"`
	NATS nats.Config `envconfig:"nats"`
}

func New() (Config, error) {
	var cfg Config
	err := envconfig.Process("lendo", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

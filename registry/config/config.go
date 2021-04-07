package config

import (
	"github.com/ivanovaleksey/lendo/registry/connectors/bank"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Bank bank.Config `envconfig:"bank"`
}

func New() (Config, error) {
	var cfg Config
	err := envconfig.Process("lendo", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

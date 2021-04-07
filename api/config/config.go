package config

import (
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr string    `required:"true"`
	DB   db.Config `envconfig:"db"`
}

func New() (Config, error) {
	var cfg Config
	err := envconfig.Process("lendo", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

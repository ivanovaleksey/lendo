package consumer

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/models"
)

type Consumer struct {
	registrator Registrator
}

type Registrator interface {
	Register(ctx context.Context, application models.Application) error
}

func New() *Consumer {
	c := &Consumer{}
	return c
}

func (c *Consumer) Run(ctx context.Context) error {
	// todo: implement
	return nil
}

func (c *Consumer) Close() error {
	// todo: implement
	return nil
}

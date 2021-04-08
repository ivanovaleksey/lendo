package nats

import (
	"context"
	"github.com/nats-io/nats.go"
)

type Consumer struct {
	client        *Client
	componentName string

	queue   string
	subject string
	subs    *nats.Subscription

	handler Handler
}

type Handler interface {
	Handle(msg *nats.Msg)
}

func NewConsumer(opts ...ConsumerOption) *Consumer {
	c := &Consumer{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Consumer) Run(ctx context.Context) error {
	subs, err := c.client.QueueSubscribe(c.subject, c.queue, c.handler.Handle)
	if err != nil {
		return err
	}
	c.subs = subs
	return nil
}

func (c *Consumer) Close() error {
	return c.subs.Drain()
}

func (c *Consumer) ComponentName() string {
	return c.componentName
}

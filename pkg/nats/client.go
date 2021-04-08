package nats

import (
	"github.com/nats-io/nats.go"
)

type Client struct {
	*nats.Conn
}

func New(cfg Config, opts ...nats.Option) (*Client, error) {
	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{Conn: nc}, err
}

func (c *Client) Close() error {
	c.Conn.Close()
	return nil
}

func (c *Client) ComponentName() string {
	return "nats"
}

package nats

type ConsumerOption func(consumer *Consumer)

func WithQueue(q string) ConsumerOption {
	return func(c *Consumer) {
		c.queue = q
	}
}

func WithSubject(s string) ConsumerOption {
	return func(c *Consumer) {
		c.subject = s
	}
}

func WithClient(cli *Client) ConsumerOption {
	return func(c *Consumer) {
		c.client = cli
	}
}

func WithHandler(h Handler) ConsumerOption {
	return func(c *Consumer) {
		c.handler = h
	}
}

func WithComponentName(name string) ConsumerOption {
	return func(c *Consumer) {
		c.componentName = name
	}
}

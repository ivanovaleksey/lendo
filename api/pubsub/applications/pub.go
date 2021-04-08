package applicationsPubSub

import (
	"context"
	"encoding/json"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/ivanovaleksey/lendo/pkg/nats"
)

type Pub struct {
	client *nats.Client
}

func NewPub(client *nats.Client) *Pub {
	pub := &Pub{
		client: client,
	}
	return pub
}

func (p *Pub) NewApplication(ctx context.Context, application models.Application) error {
	const subject = "applications.new"

	data, err := json.Marshal(application)
	if err != nil {
		return err
	}

	return p.client.Publish(subject, data)
}

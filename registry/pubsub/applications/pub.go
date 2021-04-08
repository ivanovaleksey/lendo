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

func (p *Pub) ApplicationStatusChanged(ctx context.Context, change models.StatusChange) error {
	const subject = "applications.changed"

	data, err := json.Marshal(change)
	if err != nil {
		return err
	}
	return p.client.Publish(subject, data)
}

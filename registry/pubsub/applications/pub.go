package applicationsPubSub

import (
	"context"
	"encoding/json"
	"github.com/ivanovaleksey/lendo/pkg/models"
)

type Pub struct {
	client PubClient
}

type PubClient interface {
	Publish(subj string, data []byte) error
}

func NewPub(client PubClient) *Pub {
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

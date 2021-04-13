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

func (p *Pub) NewApplication(ctx context.Context, application models.Application) error {
	const subject = "applications.new"

	data, err := json.Marshal(application)
	if err != nil {
		return err
	}

	return p.client.Publish(subject, data)
}

package poller

import "github.com/ivanovaleksey/lendo/registry/poller/handlers"

type Option func(*Poller)

func WithBank(b handlers.Bank) Option {
	return func(p *Poller) {
		p.bank = b
	}
}

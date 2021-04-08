package poller

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/ivanovaleksey/lendo/registry/poller/handlers"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	numWorkers = 2
)

type Poller struct {
	bank     handlers.Bank
	notifier handlers.Notifier

	db         *db.DB
	numWorkers int

	workersWg     sync.WaitGroup
	workersCancel context.CancelFunc
}

func New(bank handlers.Bank, db *db.DB, notifier handlers.Notifier) *Poller {
	p := &Poller{
		bank:       bank,
		db:         db,
		notifier:   notifier,
		numWorkers: numWorkers,
	}
	return p
}

func (p *Poller) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	p.workersCancel = cancel

	for i := 0; i < p.numWorkers; i++ {
		p.workersWg.Add(1)
		go func(ctx context.Context, id int) {
			defer p.workersWg.Done()

			worker := p.newWorker(id + 1)
			worker.Run(ctx)
		}(ctx, i)
	}

	return nil
}

func (p *Poller) newWorker(id int) *worker {
	logger := log.WithFields(log.Fields{
		"component": "worker",
		"id":        id,
	})
	w := &worker{
		db:     p.db,
		logger: logger,
		handlers: map[models.JobStatus]JobHandler{
			models.JobStatusNew:     handlers.NewNewJobHandler(p.bank, p.notifier),
			models.JobStatusPending: handlers.NewPendingJobHandler(p.bank, p.notifier),
		},
	}
	return w
}

func (p *Poller) Close() error {
	p.workersCancel()
	p.workersWg.Wait()
	return nil
}

func (p *Poller) ComponentName() string {
	return "poller"
}

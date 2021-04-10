package poller

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/registry/models"
	"github.com/ivanovaleksey/lendo/registry/poller/handlers"
	"github.com/ivanovaleksey/lendo/registry/poller/worker"
	"sync"
	"time"
)

const (
	defaultNumWorkers = 2
	tickerDuration    = 10 * time.Second
)

type Poller struct {
	bank          handlers.Bank
	repo          handlers.Repo
	notifier      handlers.Notifier
	workerFactory WorkerFactory
	tickerFactory TickerFactory

	db         *db.DB
	numWorkers int

	workersWg     sync.WaitGroup
	workersCancel context.CancelFunc
}

func New(opts ...Option) *Poller {
	p := &Poller{
		numWorkers:    defaultNumWorkers,
		workerFactory: stdWorkerFactory{},
		tickerFactory: stdTickerFactory{duration: tickerDuration},
	}
	for _, opt := range opts {
		opt(p)
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

			w := p.newWorker(id + 1)
			w.Run(ctx)
		}(ctx, i)
	}

	return nil
}

func (p *Poller) newWorker(id int) Worker {
	opts := []worker.Option{
		worker.WithID(id),
		worker.WithTxFactory(db.NewTxFactory(p.db)),
		worker.WithTicker(p.tickerFactory.NewTicker()),
		worker.WithHandler(models.JobStatusNew, handlers.NewNewJobHandler(p.bank, p.repo, p.notifier)),
		worker.WithHandler(models.JobStatusPending, handlers.NewPendingJobHandler(p.bank, p.repo, p.notifier)),
	}
	w := p.workerFactory.NewWorker(opts...)
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

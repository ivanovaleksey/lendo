package poller

import (
	"context"
	"sync"
)

const (
	numWorkers = 2
)

type Poller struct {
	numWorkers int

	workersWg     sync.WaitGroup
	workersCtx    context.Context
	workersCancel context.CancelFunc
}

func New() *Poller {
	p := &Poller{
		numWorkers: numWorkers,
	}
	return p
}

// poll db with
// LIMIT numworkers
// FOR UPDATE SKIP LOCKED
// and fan out tasks via channel

func (p *Poller) Run(ctx context.Context) error {
	p.workersCtx, p.workersCancel = context.WithCancel(ctx)

	for i := 0; i < p.numWorkers; i++ {
		p.workersWg.Add(1)
		go func(ctx context.Context) {
			defer p.workersWg.Done()

			worker := newWorker()
			// todo: handle error
			worker.Run(ctx)
		}(p.workersCtx)
	}

	return nil
}

func (p *Poller) Close() error {
	p.workersCancel()
	p.workersWg.Wait()
	return nil
}

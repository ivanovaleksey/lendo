package poller

import (
	"context"
	"github.com/ivanovaleksey/lendo/registry/poller/worker"
)

type WorkerFactory interface {
	NewWorker(...worker.Option) Worker
}

type Worker interface {
	Run(context.Context) error
}

type stdWorkerFactory struct{}

func (f stdWorkerFactory) NewWorker(opts ...worker.Option) Worker {
	return worker.New(opts...)
}

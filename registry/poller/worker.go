package poller

import "context"

type worker struct {
}

func newWorker() *worker {
	w := &worker{}
	return w
}

func (w *worker) Run(ctx context.Context) error {
	// todo: implement
	return nil
}

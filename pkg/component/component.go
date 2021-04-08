package component

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/closer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	CloseDelay   = 3 * time.Second
	CloseTimeout = 5 * time.Second
)

type Component interface {
	Run(context.Context) error
	Closer
}

type Closer interface {
	ComponentName() string
	Close() error
}

func Run(ctx context.Context, cmp Component) func() error {
	logger := log.WithField("component", cmp.ComponentName())

	go func() {
		logger.Debug("running")
		if err := cmp.Run(ctx); err != nil {
			logger.Errorf("error: %v", err)
			return
		}
	}()

	return func() error {
		return Close(cmp, 0)
	}
}

func Close(cmp Closer, delay time.Duration) error {
	logger := log.WithField("component", cmp.ComponentName())

	if delay > 0 {
		logger.Debug("waiting for close delay")
		time.Sleep(delay)
	}

	logger.Debug("closing")
	cl := closer.NewTimeoutCloser(cmp, CloseTimeout)
	if err := cl.Close(); err != nil {
		logger.Errorf("close error: %v", err)
		return errors.WithMessage(err, "can't close " + cmp.ComponentName())
	}
	logger.Debug("closed")
	return nil
}

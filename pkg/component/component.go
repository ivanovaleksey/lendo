package component

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/closer"
	log "github.com/sirupsen/logrus"
	"time"
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
		return Close(cmp)
	}
}

func Close(cmp Closer) error {
	const closeTimeout = 5 * time.Second
	logger := log.WithField("component", cmp.ComponentName())

	logger.Debug("closing")
	cl := closer.NewTimeoutCloser(cmp, closeTimeout)
	if err := cl.Close(); err != nil {
		logger.Errorf("close error: %v", err)
		return err
	}
	logger.Debug("closed")
	return nil
}

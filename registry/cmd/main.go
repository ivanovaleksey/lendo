package main

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/closer"
	"github.com/ivanovaleksey/lendo/registry/config"
	"github.com/ivanovaleksey/lendo/registry/consumer"
	"github.com/ivanovaleksey/lendo/registry/poller"
	log "github.com/sirupsen/logrus"
	"syscall"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)

	// todo: should cancel on close?
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	runApps(ctx, cfg)
}

func runApps(ctx context.Context, cfg config.Config) {
	appCloser := closer.New(syscall.SIGTERM, syscall.SIGINT)

	appCloser.Add(runApp(ctx, "consumer", consumer.New()))
	appCloser.Add(runApp(ctx, "poller", poller.New()))

	appCloser.Wait()
}

type App interface {
	Run(context.Context) error
	Close() error
}

func runApp(ctx context.Context, name string, app App) func() error {
	const closeTimeout = 5 * time.Second

	logger := log.WithField("component", name)

	go func() {
		logger.Debug("running")
		if err := app.Run(ctx); err != nil {
			logger.Errorf("error: %v", err)
			return
		}
	}()

	return func() error {
		logger.Debug("closing")
		cl := closer.NewTimeoutCloser(app, closeTimeout)
		if err := cl.Close(); err != nil {
			logger.Errorf("close error: %v", err)
			return err
		}
		logger.Debug("closed")
		return nil
	}
}

package main

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/closer"
	"github.com/ivanovaleksey/lendo/pkg/component"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/nats"
	"github.com/ivanovaleksey/lendo/registry/bank"
	"github.com/ivanovaleksey/lendo/registry/config"
	"github.com/ivanovaleksey/lendo/registry/poller"
	"github.com/ivanovaleksey/lendo/registry/pubsub/applications"
	"github.com/ivanovaleksey/lendo/registry/repos/jobs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"syscall"
)

func main() {
	log.SetLevel(log.DebugLevel)

	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := runApps(ctx, cfg); err != nil {
		log.Error(err)
	}
}

func runApps(ctx context.Context, cfg config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	db, err := db.New(cfg.DB)
	if err != nil {
		return errors.Wrap(err, "can't create db")
	}

	natsClient, err := nats.New(cfg.NATS)
	if err != nil {
		return errors.Wrap(err, "can't create nats client")
	}

	appCloser := closer.New(syscall.SIGTERM, syscall.SIGINT)
	appCloser.Add(func() error {
		cancel()
		return nil
	})

	{
		repo := jobsRepo.New(db)
		handler := applicationsPubSub.NewNewApplicationHandler(repo)

		opts := []nats.ConsumerOption{
			nats.WithClient(natsClient),
			nats.WithQueue("registry"),
			nats.WithSubject("applications.new"),
			nats.WithHandler(handler),
			nats.WithComponentName("consumer.applications.new"),
		}
		closure := component.Run(ctx, nats.NewConsumer(opts...))
		appCloser.Add(closure)
	}

	{
		bankClient := bank.NewClient(cfg.Bank)
		repo := jobsRepo.New(db)
		pub := applicationsPubSub.NewPub(natsClient)

		opts := []poller.Option{
			poller.WithDB(db),
			poller.WithBank(bankClient),
			poller.WithRepo(repo),
			poller.WithNotifier(pub),
		}
		closure := component.Run(ctx, poller.New(opts...))
		appCloser.Add(closure)
	}

	appCloser.Add(func() error {
		return component.Close(natsClient, component.CloseDelay)
	})
	appCloser.Add(func() error {
		return component.Close(db, component.CloseDelay)
	})

	appCloser.Wait()
	return nil
}

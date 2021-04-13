//go:generate swagger generate spec -o ../docs/swagger.json

package main

import (
	"context"
	"github.com/ivanovaleksey/lendo/api/app"
	"github.com/ivanovaleksey/lendo/api/config"
	"github.com/ivanovaleksey/lendo/api/pubsub/applications"
	"github.com/ivanovaleksey/lendo/api/repos/applications"
	"github.com/ivanovaleksey/lendo/api/services/applications"
	"github.com/ivanovaleksey/lendo/pkg/closer"
	"github.com/ivanovaleksey/lendo/pkg/component"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/nats"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"syscall"
	"time"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
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

	repo := applicationsRepo.New(db)

	var opts []app.Option
	{
		pub := applicationsPubSub.NewPub(natsClient)
		srv := applicationsSrv.New(repo, pub)
		opts = append(opts, app.WithApplicationsSrv(srv))
	}

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      app.New(cfg, opts...),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	appCloser := closer.New(syscall.SIGTERM, syscall.SIGINT)
	appCloser.Add(func() error {
		cancel()
		return nil
	})

	{
		handler := applicationsPubSub.NewApplicationStatusChangedHandler(repo)

		opts := []nats.ConsumerOption{
			nats.WithClient(natsClient),
			nats.WithQueue("api"),
			nats.WithSubject("applications.changed"),
			nats.WithHandler(handler),
			nats.WithComponentName("consumer.applications.changed"),
		}
		closure := component.Run(ctx, nats.NewConsumer(opts...))
		appCloser.Add(closure)
	}

	appCloser.Add(func() error {
		return closeSrv(&srv)
	})
	appCloser.Add(func() error {
		return component.Close(natsClient, component.CloseDelay)
	})
	appCloser.Add(func() error {
		return component.Close(db, component.CloseDelay)
	})

	go func() {
		log.Debugf("starting server on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("server error: %v", err)
			appCloser.CloseAll()
		}
	}()

	appCloser.Wait()
	return nil
}

func closeSrv(srv *http.Server) error {
	const (
		gracefulDelay   = 3 * time.Second
		gracefulTimeout = 5 * time.Second
	)

	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	// waiting for k8s to stop traffic
	log.Info("waiting for graceful delay")
	time.Sleep(gracefulDelay)

	log.Info("shutting down")
	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown error")
	}
	log.Info("shutdown gracefully")
	return nil
}

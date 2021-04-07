package main

import (
	"context"
	"github.com/ivanovaleksey/lendo/api/app"
	"github.com/ivanovaleksey/lendo/api/config"
	"github.com/ivanovaleksey/lendo/pkg/closer"
	"github.com/ivanovaleksey/lendo/pkg/db"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := runApp(ctx, cfg); err != nil {
		log.Error(err)
	}
}

func runApp(ctx context.Context, cfg config.Config) error {
	db, err := db.New(cfg.DB)
	if err != nil {
		return errors.Wrap(err, "can't create db")
	}

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      app.New(cfg, db),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	appCloser := closer.New(syscall.SIGTERM, syscall.SIGINT)
	appCloser.Add(func() error {
		return closeApp(ctx, &srv)
	})
	appCloser.Add(func() error {
		const closeTimeout = 3 * time.Second
		logger := log.WithField("component", "db")

		logger.Debug("closing")
		cl := closer.NewTimeoutCloser(db, closeTimeout)
		if err := cl.Close(); err != nil {
			logger.Errorf("close error: %v", err)
			return err
		}
		logger.Debug("closed")
		return nil
	})

	go func() {
		log.Debugf("starting server on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("server error: %v", err)
			appCloser.CloseAll()
		}
	}()

	appCloser.Wait()
	// todo: is it possible to return close error?
	return nil
}

func closeApp(ctx context.Context, srv *http.Server) error {
	const (
		gracefulDelay   = 3 * time.Second
		gracefulTimeout = 5 * time.Second
	)

	ctx, cancel := context.WithTimeout(ctx, gracefulTimeout)
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

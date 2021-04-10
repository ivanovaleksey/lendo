package poller_test

import (
	"context"
	"github.com/ivanovaleksey/lendo/pkg/db"
	"github.com/ivanovaleksey/lendo/pkg/test"
	"github.com/ivanovaleksey/lendo/pkg/ticker/mocks"
	"github.com/ivanovaleksey/lendo/registry/config"
	"github.com/ivanovaleksey/lendo/registry/poller"
	mockHandlers "github.com/ivanovaleksey/lendo/registry/poller/handlers/mocks"
	"github.com/ivanovaleksey/lendo/registry/poller/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPoller_Run(t *testing.T) {
	t.Run("should spawn workers", func(t *testing.T) {
		const numWorkers = 3

		opts := []poller.Option{
			poller.WithNumWorkers(numWorkers),
		}

		fx := newFixture(t, opts...)
		defer fx.Finish()

		tick := &mockTicker.Ticker{}
		fx.tickerFactory.On("NewTicker").Return(tick).Times(numWorkers)

		for i := 0; i < numWorkers; i++ {
			wrk := &mocks.Worker{}
			wrk.On("Run", mock.AnythingOfType("*context.cancelCtx")).Return(nil).Once()
			fx.workers = append(fx.workers, wrk)

			fx.workerFactory.On("NewWorker", mock.AnythingOfType("[]worker.Option")).Return(wrk).Once()
		}

		err := fx.poller.Run(fx.ctx)

		require.NoError(t, err)
		fx.Wait()
	})
}

type fixture struct {
	t   *testing.T
	ctx context.Context
	db  *db.DB

	bank          *mockHandlers.Bank
	notifier      *mockHandlers.Notifier
	tickerFactory *mocks.TickerFactory
	workerFactory *mocks.WorkerFactory
	workers       []*mocks.Worker

	poller *poller.Poller
}

func newFixture(t *testing.T, opts ...poller.Option) *fixture {
	test.LoadRegistryEnv(t)

	cfg, err := config.New()
	require.NoError(t, err)

	fx := &fixture{
		t:   t,
		ctx: context.Background(),
		db:  db.NewTestDB(t, cfg.DB),

		tickerFactory: &mocks.TickerFactory{},
		workerFactory: &mocks.WorkerFactory{},
		bank:          &mockHandlers.Bank{},
		notifier:      &mockHandlers.Notifier{},
	}

	baseOpts := []poller.Option{
		poller.WithDB(fx.db),
		poller.WithBank(fx.bank),
		poller.WithNotifier(fx.notifier),
		poller.WithNumWorkers(1),
		poller.WithWorkerFactory(fx.workerFactory),
		poller.WithTickerFactory(fx.tickerFactory),
	}
	opts = append(baseOpts, opts...)
	fx.poller = poller.New(opts...)

	return fx
}

func (fx *fixture) Finish() {
	require.NoError(fx.t, fx.poller.Close())
	require.True(fx.t, fx.tickerFactory.AssertExpectations(fx.t))
	require.True(fx.t, fx.workerFactory.AssertExpectations(fx.t))
	for _, w := range fx.workers {
		w.AssertExpectations(fx.t)
	}
}

func (fx *fixture) Wait() {
	done := make(chan struct{})
	time.AfterFunc(1*time.Second, func() {
		defer close(done)
		require.NoError(fx.t, fx.poller.Close())
	})
	<-done
}

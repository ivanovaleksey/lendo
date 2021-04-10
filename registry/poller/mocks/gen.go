//go:generate mockery --dir .. --output . --name Worker --filename worker.mock.go
//go:generate mockery --dir .. --output . --name WorkerFactory --unroll-variadic=false --filename worker_factory.mock.go
//go:generate mockery --dir .. --output . --name TickerFactory --filename ticker_factory.mock.go

package mocks

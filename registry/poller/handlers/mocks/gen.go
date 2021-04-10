//go:generate mockery --dir .. --output . --outpkg mockHandlers --name Bank --filename bank.mock.go
//go:generate mockery --dir .. --output . --outpkg mockHandlers --name Repo --filename repo.mock.go
//go:generate mockery --dir .. --output . --outpkg mockHandlers --name Notifier --filename notifier.mock.go

package mockHandlers

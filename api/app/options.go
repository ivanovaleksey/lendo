package app

type Option func(*API)

func WithApplicationsSrv(srv ApplicationsService) Option {
	return func(api *API) {
		api.applicationsSrv = srv
	}
}

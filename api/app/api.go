package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ivanovaleksey/lendo/api/config"
	"net/http"
)

type API struct {
	cfg       config.Config
	router    chi.Router
	validator *validator.Validate

	applicationsSrv ApplicationsService
}

func New(cfg config.Config, opts ...Option) *API {
	app := &API{
		cfg:       cfg,
		validator: validator.New(),
	}
	for _, opt := range opts {
		opt(app)
	}
	app.initRouter()
	return app
}

func (api *API) initRouter() {
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Route("/applications", func(r chi.Router) {
			r.Get("/", api.GetApplications())
			r.Get("/{id}", api.GetApplication())
			r.Post("/", api.CreateApplication())
		})
	})
	api.router = router
}

func (api *API) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	api.router.ServeHTTP(writer, request)
}

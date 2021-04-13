package app

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	apiModels "github.com/ivanovaleksey/lendo/api/models"
	"github.com/ivanovaleksey/lendo/api/responses"
	"github.com/ivanovaleksey/lendo/api/services/applications"
	"github.com/ivanovaleksey/lendo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

type ApplicationsService interface {
	GetList(ctx context.Context, params applicationsSrv.GetListParams) ([]models.Application, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Application, error)
	Create(ctx context.Context, item models.NewApplication) (uuid.UUID, error)
}

// swagger:parameters getApplications
type GetApplicationsParams struct {
	apiModels.PaginationParams
	// Application status
	// in: query
	Status string `json:"status"`
}

// swagger:route GET /applications getApplications
//
// Lists applications filtered by status.
//
//     Responses:
//       default: errorResponse
//       200: getApplicationsResponse
func (api *API) GetApplications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		params := GetApplicationsParams{
			Status: r.URL.Query().Get("status"),
		}
		if value := r.URL.Query().Get("offset"); value != "" {
			offset, err := strconv.Atoi(value)
			if err != nil {
				render.Render(w, r, responses.ErrBadRequest(err))
				return
			}
			params.Offset = offset
		}
		if value := r.URL.Query().Get("limit"); value != "" {
			limit, err := strconv.Atoi(value)
			if err != nil {
				render.Render(w, r, responses.ErrBadRequest(err))
				return
			}
			params.Limit = limit
		}

		applications, total, err := api.applicationsSrv.GetList(ctx, applicationsSrv.GetListParams{
			PaginationParams: params.PaginationParams,
			Status:           params.Status,
		})
		if err != nil {
			render.Render(w, r, responses.ErrInternal(err))
			return
		}

		resp := responses.GetApplicationsResponse{
			Items: applications,
			Total: total,
		}
		render.Render(w, r, resp)
		return
	}
}

// swagger:parameters getApplication
type GetApplicationParams struct {
	// required: true
	// in: path
	ID uuid.UUID `json:"id"`
}

// swagger:route GET /applications/{id} getApplication
//
// Get application by ID.
//
//     Responses:
//       default: errorResponse
//       200: getApplicationResponse
func (api *API) GetApplication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.FromString(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.ErrBadRequest(err))
			return
		}

		application, err := api.applicationsSrv.GetByID(ctx, id)

		resp := responses.GetApplicationResponse{Application: application}
		render.Render(w, r, resp)
		return
	}
}

// swagger:parameters createApplication
type CreateApplicationParams struct {
	// in: body
	Body models.NewApplication
}

// swagger:route POST /applications createApplication
//
// Create application.
//
//     Responses:
//       default: errorResponse
//       200: createApplicationResponse
func (api *API) CreateApplication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		ctx := r.Context()

		var params CreateApplicationParams
		err := json.NewDecoder(r.Body).Decode(&params.Body)
		if err != nil {
			render.Render(w, r, responses.ErrBadRequest(err))
			return
		}

		err = api.validator.Struct(params)
		if err != nil {
			render.Render(w, r, responses.ErrBadRequest(err))
			return
		}

		id, err := api.applicationsSrv.Create(ctx, params.Body)

		resp := responses.CreateApplicationResponse{ID: id}
		render.Render(w, r, resp)
		return
	}
}

package app

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	applicationsRepo "github.com/ivanovaleksey/lendo/api/repos/applications"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type GetApplicationsResponse struct {
	Items []models.Application `json:"items"`
	Total int                  `json:"total"`
}

func (api *API) GetApplications() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		var params applicationsRepo.GetListParams
		if value := request.URL.Query().Get("offset"); value != "" {
			offset, err := strconv.Atoi(value)
			if err != nil {
				log.Error(errors.Wrap(err, "can't parse offset parameter"))
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			params.Offset = offset
		}
		if value := request.URL.Query().Get("limit"); value != "" {
			limit, err := strconv.Atoi(value)
			if err != nil {
				log.Error(errors.Wrap(err, "can't parse limit parameter"))
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			params.Limit = limit
		}
		status := request.URL.Query().Get("status")
		if status == "" {
			log.Error("empty status")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		params.Status = status

		applications, total, err := api.applicationsRepo.GetList(ctx, params)
		if err != nil {
			log.Error(errors.Wrap(err, "can't get applications"))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := GetApplicationsResponse{
			Items: applications,
			Total: total,
		}
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Error(errors.Wrap(err, "can't encode response"))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

type GetApplicationResponse struct {
	models.Application
}

func (api *API) GetApplication() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		id, err := uuid.FromString(chi.URLParam(request, "id"))
		if err != nil {
			log.Error(errors.Wrap(err, "invalid id"))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		application, err := api.applicationsRepo.GetByID(ctx, id)
		err = json.NewEncoder(writer).Encode(application)
		if err != nil {
			log.Error(errors.Wrap(err, "can't encode response"))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

type CreateApplicationResponse struct {
	ID uuid.UUID `json:"id"`
}

func (api *API) CreateApplication() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		ctx := request.Context()

		var application models.NewApplication
		err := json.NewDecoder(request.Body).Decode(&application)
		if err != nil {
			log.Error(errors.Wrap(err, "can't decode body"))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = api.validator.Struct(application)
		if err != nil {
			// todo: show errors
			log.Error(errors.Wrap(err, "can't validate"))
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := api.applicationsRepo.Create(ctx, application)
		resp := CreateApplicationResponse{ID: id}
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Error(errors.Wrap(err, "can't encode response"))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

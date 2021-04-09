package bank

import (
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit"
	"github.com/ivanovaleksey/lendo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImpl_CreateApplication(t *testing.T) {
	application := models.Application{
		ID:     uuid.NewV4(),
		Status: models.ApplicationStatus(gofakeit.Word()),
		NewApplication: models.NewApplication{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		},
	}

	t.Run("when application already exists", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		errMsg := gofakeit.Sentence(3)

		serverMock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "POST", r.Method)
			require.Equal(t, "/api/applications", r.URL.String())

			var body createApplicationBody
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.Equal(t, application.ID, body.ID)
			require.Equal(t, application.FirstName, body.FirstName)
			require.Equal(t, application.LastName, body.LastName)

			resp := `{
				"error": "` + errMsg + `"
			}`
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer serverMock.Close()

		status, err := fx.client.CreateApplication(fx.ctx, application)

		expectedErr := Error{
			Code:    http.StatusBadRequest,
			Message: errMsg,
		}
		require.Equal(t, expectedErr, err)
		assert.Empty(t, status)
	})

	t.Run("when application does not exist", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		status := gofakeit.Word()

		serverMock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "POST", r.Method)
			require.Equal(t, "/api/applications", r.URL.String())

			var body createApplicationBody
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.Equal(t, application.ID, body.ID)
			require.Equal(t, application.FirstName, body.FirstName)
			require.Equal(t, application.LastName, body.LastName)

			resp := `{
				"id": "` + application.ID.String() + `",
				"first_name": "` + application.FirstName + `",
				"last_name": "` + application.LastName + `",
				"status": "` + status + `"
			}`
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer serverMock.Close()

		got, err := fx.client.CreateApplication(fx.ctx, application)

		require.NoError(t, err)
		assert.EqualValues(t, status, got)
	})

	t.Run("with unknown error", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		serverMock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "POST", r.Method)
			require.Equal(t, "/api/applications", r.URL.String())

			var body createApplicationBody
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.Equal(t, application.ID, body.ID)
			require.Equal(t, application.FirstName, body.FirstName)
			require.Equal(t, application.LastName, body.LastName)

			resp := `{}`
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer serverMock.Close()

		status, err := fx.client.CreateApplication(fx.ctx, application)

		require.Equal(t, Error{Code: http.StatusInternalServerError}, err)
		assert.Empty(t, status)
	})
}

func TestImpl_GetApplicationStatus(t *testing.T) {
	applicationID := uuid.NewV4()

	t.Run("when application_id is invalid", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		errMsg := gofakeit.Sentence(3)

		serverMock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "/api/jobs?application_id="+applicationID.String(), r.URL.String())

			resp := `{
				"error": "` + errMsg + `"
			}`
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer serverMock.Close()

		status, err := fx.client.GetApplicationStatus(fx.ctx, applicationID)

		expectedErr := Error{
			Code:    http.StatusBadRequest,
			Message: errMsg,
		}
		require.Equal(t, expectedErr, err)
		assert.Empty(t, status)
	})

	t.Run("when application does not exist", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		errMsg := gofakeit.Sentence(3)

		serverMock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "/api/jobs?application_id="+applicationID.String(), r.URL.String())

			resp := `{
				"error": "` + errMsg + `"
			}`
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer serverMock.Close()

		status, err := fx.client.GetApplicationStatus(fx.ctx, applicationID)

		expectedErr := Error{
			Code:    http.StatusNotFound,
			Message: errMsg,
		}
		require.Equal(t, expectedErr, err)
		assert.Empty(t, status)
	})

	t.Run("with unknown error", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		serverMock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "/api/jobs?application_id="+applicationID.String(), r.URL.String())

			resp := `{}`
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer serverMock.Close()

		status, err := fx.client.GetApplicationStatus(fx.ctx, applicationID)

		require.Equal(t, Error{Code: http.StatusInternalServerError}, err)
		assert.Empty(t, status)
	})

	t.Run("when everything is fine", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish()

		status := gofakeit.Word()

		serverMock := fx.newServerMock(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "/api/jobs?application_id="+applicationID.String(), r.URL.String())

			resp := `{
				"id": "` + uuid.NewV4().String() + `",
				"application_id": "` + applicationID.String() + `",
				"status": "` + status + `"
			}`
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(resp))
			require.NoError(t, err)
		})
		defer serverMock.Close()

		got, err := fx.client.GetApplicationStatus(fx.ctx, applicationID)

		require.NoError(t, err)
		assert.EqualValues(t, status, got)
	})
}

type fixture struct {
	t   *testing.T
	ctx context.Context

	client *Client
}

func newFixture(t *testing.T) *fixture {
	client := NewClient(Config{})

	fx := &fixture{
		t:      t,
		ctx:    context.Background(),
		client: client,
	}
	return fx
}

func (fx *fixture) Finish() {}

func (fx *fixture) newServerMock(fn http.HandlerFunc) *httptest.Server {
	srv := httptest.NewServer(fn)
	fx.client.cfg.URL = srv.URL
	return srv
}

package bank

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ivanovaleksey/lendo/pkg/models"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"path"
	"time"
)

type impl struct {
	cfg        Config
	httpClient *http.Client
}

func New(cfg Config) Client {
	client := &impl{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
	return client
}

func (client *impl) CreateApplication(ctx context.Context, application models.Application) (models.ApplicationStatus, error) {
	const methodPath = "/api/applications"

	reqURL, err := url.Parse(client.cfg.URL)
	if err != nil {
		return "", errors.Wrap(err, "can't parse url")
	}
	reqURL.Path = path.Join(reqURL.Path, methodPath)

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(application); err != nil {
		return "", errors.Wrap(err, "can't encode request body")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), &body)
	if err != nil {
		return "", errors.Wrap(err, "can't create request")
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't do request")
	}
	defer resp.Body.Close()

	var respBody struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", errors.Wrap(err, "can't decode response body")
	}

	return models.ApplicationStatus(respBody.Status), nil
}

func (client *impl) GetApplicationStatus(ctx context.Context, id uuid.UUID) (models.ApplicationStatus, error) {
	const methodPath = "/api/jobs"

	reqURL, err := url.Parse(client.cfg.URL)
	if err != nil {
		return "", errors.Wrap(err, "can't parse url")
	}
	reqURL.Path = path.Join(reqURL.Path, methodPath)

	params := url.Values{
		"application_id": []string{id.String()},
	}
	reqURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return "", errors.Wrap(err, "can't create request")
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't do request")
	}
	defer resp.Body.Close()

	var respBody struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", errors.Wrap(err, "can't decode response body")
	}

	return models.ApplicationStatus(respBody.Status), nil
}

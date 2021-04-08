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

// Client interacts with a bank system.
type Client struct {
	cfg        Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	client := &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
	return client
}

func (client *Client) CreateApplication(ctx context.Context, application models.Application) (models.ApplicationStatus, error) {
	type requestBody struct {
		ID        uuid.UUID `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
	}

	const methodPath = "/api/applications"

	reqURL, err := url.Parse(client.cfg.URL)
	if err != nil {
		return "", errors.Wrap(err, "can't parse url")
	}
	reqURL.Path = path.Join(reqURL.Path, methodPath)

	var reqBody bytes.Buffer
	var reqPayload = requestBody{
		ID:        application.ID,
		FirstName: application.FirstName,
		LastName:  application.LastName,
	}
	if err := json.NewEncoder(&reqBody).Encode(reqPayload); err != nil {
		return "", errors.Wrap(err, "can't encode request body")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), &reqBody)
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

func (client *Client) GetApplicationStatus(ctx context.Context, id uuid.UUID) (models.ApplicationStatus, error) {
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

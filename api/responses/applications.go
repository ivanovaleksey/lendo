package responses

import (
	"github.com/ivanovaleksey/lendo/pkg/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

// Get applications list response
// swagger:response getApplicationsResponse
type GetApplicationsResponse_ struct {
	// in: body
	Body GetApplicationsResponse
}

type GetApplicationsResponse struct {
	Items []models.Application `json:"items"`
	Total int                  `json:"total"`
}

func (GetApplicationsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Get application by ID response
// swagger:response getApplicationResponse
type GetApplicationResponse_ struct {
	// in: body
	Body GetApplicationResponse
}

type GetApplicationResponse struct {
	models.Application
}

func (GetApplicationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Create application response
// swagger:response createApplicationResponse
type CreateApplicationResponse_ struct {
	// in: body
	Body CreateApplicationResponse
}

type CreateApplicationResponse struct {
	ID uuid.UUID `json:"id"`
}

func (CreateApplicationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

package models

import (
	"github.com/ivanovaleksey/lendo/pkg/models"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID          uuid.UUID          `json:"id"`
	Application models.Application `json:"application"`
	Status      JobStatus          `json:"status"`
}

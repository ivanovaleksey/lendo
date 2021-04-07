package models

import "github.com/satori/go.uuid"

type NewApplication struct {
	FirstName string `json:"first_name" db:"first_name" validate:"required"`
	LastName  string `json:"last_name" db:"last_name" validate:"required"`
}

type Application struct {
	NewApplication
	ID     uuid.UUID         `json:"id"`
	Status ApplicationStatus `json:"status,omitempty"`
}

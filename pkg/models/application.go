package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
)

type NewApplication struct {
	FirstName string `json:"first_name" db:"first_name" validate:"required"`
	LastName  string `json:"last_name" db:"last_name" validate:"required"`
}

type Application struct {
	NewApplication
	ID     uuid.UUID         `json:"id"`
	Status ApplicationStatus `json:"status"`
}

func (a Application) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Application) Scan(src interface{}) error {
	data, ok := src.([]byte)
	if !ok {
		return errors.New("expected []byte")
	}
	return json.Unmarshal(data, &a)
}

type StatusChange struct {
	ID     uuid.UUID         `json:"id"`
	Status ApplicationStatus `json:"status"`
}

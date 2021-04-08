package models

type ApplicationStatus string

const (
	ApplicationStatusNew       = "new"
	ApplicationStatusPending   = "pending"
	ApplicationStatusCompleted = "completed"
	ApplicationStatusRejected  = "rejected"
)

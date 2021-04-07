package models

type ApplicationStatus string

const (
	ApplicationStatusUnknown   = ""
	ApplicationStatusPending   = "pending"
	ApplicationStatusCompleted = "completed"
	ApplicationStatusRejected  = "rejected"
)

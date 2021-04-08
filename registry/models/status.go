package models

type JobStatus string

const (
	JobStatusNew     = "new"
	JobStatusPending = "pending"
	JobStatusDone    = "done"
)

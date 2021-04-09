package models

type JobStatus string

const (
	JobStatusNew     JobStatus = "new"
	JobStatusPending JobStatus = "pending"
	JobStatusDone    JobStatus = "done"
)

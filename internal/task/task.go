package task

import "time"

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusClaimed   Status = "CLAIMED"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

type Task struct {
	ID        string
	Payload   string
	Status    Status
	AgentID   string
	Attempt   int
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
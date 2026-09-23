package store

import "github.com/mthatipamula/go-agent-control-plane/internal/task"

type TaskRepository interface {
	Create(t task.Task) error
	Get(id string) (task.Task, error)
	UpdateWithFencing(t task.Task, expectedVersion int64, expectedFencingToken int64) error
	ListPending() ([]task.Task, error)
	List() ([]task.Task, error)
}

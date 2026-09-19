package agent

import (
	"errors"

	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

type Agent struct {
	ID    string
	store *store.TaskStore
}

func NewAgent(id string, store *store.TaskStore) *Agent {
	return &Agent{
		ID:    id,
		store: store,
	}
}

func (a *Agent) Claim(taskID string) error {
	t, err := a.store.Get(taskID)
	if err != nil {
		return err
	}

	if t.Status != task.StatusPending {
		return errors.New("task is not pending")
	}

	t.Status = task.StatusClaimed
	t.AgentID = a.ID

	return a.store.Update(t, t.Version)
}

func (a *Agent) Start(taskID string) error {
	t, err := a.store.Get(taskID)
	if err != nil {
		return err
	}

	if t.Status != task.StatusClaimed {
		return errors.New("task is not claimed")
	}

	if t.AgentID != a.ID {
		return errors.New("task is claimed by another agent")
	}

	t.Status = task.StatusRunning

	return a.store.Update(t, t.Version)
}

func (a *Agent) Complete(taskID string) error {
	t, err := a.store.Get(taskID)
	if err != nil {
		return err
	}

	if t.Status != task.StatusRunning {
		return errors.New("task is not running")
	}

	if t.AgentID != a.ID {
		return errors.New("task is owned by another agent")
	}

	t.Status = task.StatusCompleted

	return a.store.Update(t, t.Version)
}

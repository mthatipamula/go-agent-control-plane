package agent

import (
	"errors"
	"time"

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

	now := time.Now()
	leaseExpiresAt := now.Add(LeaseDuration)

	t.Status = task.StatusClaimed
	t.AgentID = a.ID
	t.LeaseExpiresAt = &leaseExpiresAt
	t.FencingToken++

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

func (a *Agent) Renew(taskID string) error {
	t, err := a.store.Get(taskID)
	if err != nil {
		return err
	}

	if t.AgentID != a.ID {
		return errors.New("task is owned by another agent")
	}

	if t.LeaseExpiresAt == nil {
		return errors.New("task has no lease")
	}

	leaseExpiresAt := time.Now().Add(LeaseDuration)
	t.LeaseExpiresAt = &leaseExpiresAt

	return a.store.Update(t, t.Version)
}

func (a *Agent) Recover(taskID string) error {
	t, err := a.store.Get(taskID)
	if err != nil {
		return err
	}

	if !t.IsLeaseExpired(time.Now()) {
		return errors.New("task lease has not expired")
	}

	t.AgentID = a.ID
	t.Status = task.StatusClaimed
	t.FencingToken++

	leaseExpiresAt := time.Now().Add(LeaseDuration)
	t.LeaseExpiresAt = &leaseExpiresAt

	return a.store.Update(t, t.Version)
}

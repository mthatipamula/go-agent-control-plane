package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

type Agent struct {
	ID       string
	store    store.TaskRepository
	executor task.Executor
}

func NewAgent(id string, store store.TaskRepository, executors ...task.Executor) *Agent {
	var executor task.Executor

	if len(executors) > 0 {
		executor = executors[0]
	}

	return &Agent{
		ID:       id,
		store:    store,
		executor: executor,
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

	return a.store.UpdateWithFencing(t, t.Version, t.FencingToken-1)
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

	return a.store.UpdateWithFencing(t, t.Version, t.FencingToken)
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

	return a.store.UpdateWithFencing(t, t.Version, t.FencingToken)
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

	return a.store.UpdateWithFencing(t, t.Version, t.FencingToken)
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

	return a.store.UpdateWithFencing(t, t.Version, t.FencingToken-1)
}

func (a *Agent) Execute(taskID string) error {
	t, err := a.store.Get(taskID)
	if err != nil {
		return err
	}

	if t.Status != task.StatusRunning {
		return fmt.Errorf("task %s is not running", taskID)
	}

	if t.AgentID != a.ID {
		return fmt.Errorf("task %s is not owned by agent %s", taskID, a.ID)
	}

	if a.executor == nil {
		return fmt.Errorf("agent %s has no executor", a.ID)
	}

	if err := a.executor.Execute(t.Payload); err != nil {
		t.Status = task.StatusFailed
		t.LeaseExpiresAt = nil

		if err := a.store.UpdateWithFencing(
			t,
			t.Version,
			t.FencingToken,
		); err != nil {
			return err
		}

		return err
	}

	t.Status = task.StatusCompleted
	t.LeaseExpiresAt = nil

	return a.store.UpdateWithFencing(
		t,
		t.Version,
		t.FencingToken,
	)
}

func (a *Agent) RunOnce() error {
	pending, err := a.store.ListPending()
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		return nil
	}

	t := pending[0]

	if err := a.Claim(t.ID); err != nil {
		return err
	}

	if err := a.Start(t.ID); err != nil {
		return err
	}

	return a.Execute(t.ID)
}

func (a *Agent) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := a.RunOnce(); err != nil {
				log.Printf("agent %s: task execution error: %v", a.ID, err)
			}
		}
	}
}

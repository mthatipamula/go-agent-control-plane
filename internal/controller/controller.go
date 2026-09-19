package controller

import (
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

type Controller struct {
	store *store.TaskStore
}

func NewController(store *store.TaskStore) *Controller {
	return &Controller{
		store: store,
	}
}

func (c *Controller) Submit(payload string) (task.Task, error) {
	now := time.Now()

	t := task.Task{
		ID:        task.NewID(),
		Payload:   payload,
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.store.Create(t); err != nil {
		return task.Task{}, err
	}

	return t, nil
}
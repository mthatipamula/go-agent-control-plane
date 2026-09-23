package controller

import (
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

type Controller struct {
	store store.TaskRepository
}

func NewController(store store.TaskRepository) *Controller {
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

func (c *Controller) Get(taskID string) (task.Task, error) {
	return c.store.Get(taskID)
}

func (c *Controller) List() ([]task.Task, error) {
	return c.store.List()
}

package store

import (
	"testing"
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

func TestTaskStoreCreateAndGet(t *testing.T) {
	store := NewTaskStore()

	now := time.Now()

	input := task.Task{
		ID:        "task-1",
		Payload:   "process order",
		Status:    task.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := store.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := store.Get("task-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.ID != input.ID {
		t.Errorf("ID = %q, want %q", got.ID, input.ID)
	}

	if got.Status != task.StatusPending {
		t.Errorf("Status = %q, want %q", got.Status, task.StatusPending)
	}
}

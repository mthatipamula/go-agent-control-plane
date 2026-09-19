package store

import (
	"testing"
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

func TestTaskStoreUpdateVersionConflict(t *testing.T) {
	store := NewTaskStore()

	now := time.Now()

	input := task.Task{
		ID:        "task-1",
		Payload:   "process order",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := store.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	input.Status = task.StatusRunning

	if err := store.Update(input, 0); err != nil {
		t.Fatalf("first Update() error = %v", err)
	}

	input.Status = task.StatusCompleted

	err := store.Update(input, 0)
	if err != ErrVersionConflict {
		t.Fatalf("Update() error = %v, want %v", err, ErrVersionConflict)
	}
}
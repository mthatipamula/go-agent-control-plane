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

func TestTaskStoreUpdateWithFencingRejectsStaleToken(t *testing.T) {
	store := NewTaskStore()

	now := time.Now()

	input := task.Task{
		ID:           "task-fencing",
		Payload:      "process order",
		Status:       task.StatusClaimed,
		Version:      0,
		FencingToken: 2,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := store.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err := store.UpdateWithFencing(input, 0, 1)
	if err != ErrFencingTokenConflict {
		t.Fatalf("UpdateWithFencing() error = %v, want %v",
			err, ErrFencingTokenConflict)
	}

	got, err := store.Get("task-fencing")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.FencingToken != 2 {
		t.Errorf("FencingToken = %d, want %d", got.FencingToken, 2)
	}

	if got.Version != 0 {
		t.Errorf("Version = %d, want %d", got.Version, 0)
	}
}

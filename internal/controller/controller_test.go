package controller

import (
	"testing"

	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

func TestSubmit(t *testing.T) {
	taskStore := store.NewTaskStore()
	controller := NewController(taskStore)

	created, err := controller.Submit("process order")
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}

	if created.ID == "" {
		t.Fatal("Submit() returned an empty task ID")
	}

	if created.Payload != "process order" {
		t.Errorf("Payload = %q, want %q", created.Payload, "process order")
	}

	if created.Status != task.StatusPending {
		t.Errorf("Status = %q, want %q", created.Status, task.StatusPending)
	}

	stored, err := taskStore.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if stored.ID != created.ID {
		t.Errorf("stored ID = %q, want %q", stored.ID, created.ID)
	}
}

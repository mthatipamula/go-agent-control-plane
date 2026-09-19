package agent

import (
	"testing"
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

func TestAgentClaim(t *testing.T) {
	taskStore := store.NewTaskStore()
	agent := NewAgent("agent-1", taskStore)

	now := time.Now()

	input := task.Task{
		ID:        "task-1",
		Payload:   "process order",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := taskStore.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := agent.Claim("task-1"); err != nil {
		t.Fatalf("Claim() error = %v", err)
	}

	got, err := taskStore.Get("task-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.Status != task.StatusClaimed {
		t.Errorf("Status = %q, want %q", got.Status, task.StatusClaimed)
	}

	if got.AgentID != "agent-1" {
		t.Errorf("AgentID = %q, want %q", got.AgentID, "agent-1")
	}

	if got.Version != 1 {
		t.Errorf("Version = %d, want %d", got.Version, 1)
	}
}

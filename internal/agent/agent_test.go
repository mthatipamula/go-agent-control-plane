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

func TestAgentStart(t *testing.T) {
	taskStore := store.NewTaskStore()
	agent := NewAgent("agent-1", taskStore)

	now := time.Now()

	input := task.Task{
		ID:        "task-2",
		Payload:   "process payment",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := taskStore.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := agent.Claim("task-2"); err != nil {
		t.Fatalf("Claim() error = %v", err)
	}

	if err := agent.Start("task-2"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	got, err := taskStore.Get("task-2")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.Status != task.StatusRunning {
		t.Errorf("Status = %q, want %q", got.Status, task.StatusRunning)
	}

	if got.AgentID != "agent-1" {
		t.Errorf("AgentID = %q, want %q", got.AgentID, "agent-1")
	}

	if got.Version != 2 {
		t.Errorf("Version = %d, want %d", got.Version, 2)
	}
}

func TestAgentComplete(t *testing.T) {
	taskStore := store.NewTaskStore()
	agent := NewAgent("agent-1", taskStore)

	now := time.Now()

	input := task.Task{
		ID:        "task-3",
		Payload:   "process order",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := taskStore.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := agent.Claim("task-3"); err != nil {
		t.Fatalf("Claim() error = %v", err)
	}

	if err := agent.Start("task-3"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if err := agent.Complete("task-3"); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	got, err := taskStore.Get("task-3")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.Status != task.StatusCompleted {
		t.Errorf("Status = %q, want %q", got.Status, task.StatusCompleted)
	}

	if got.AgentID != "agent-1" {
		t.Errorf("AgentID = %q, want %q", got.AgentID, "agent-1")
	}

	if got.Version != 3 {
		t.Errorf("Version = %d, want %d", got.Version, 3)
	}
}

func TestOnlyOneAgentCanClaimTask(t *testing.T) {
	taskStore := store.NewTaskStore()

	agent1 := NewAgent("agent-1", taskStore)
	agent2 := NewAgent("agent-2", taskStore)

	now := time.Now()

	input := task.Task{
		ID:        "task-concurrent",
		Payload:   "process order",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := taskStore.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := agent1.Claim("task-concurrent"); err != nil {
		t.Fatalf("agent1 Claim() error = %v", err)
	}

	if err := agent2.Claim("task-concurrent"); err == nil {
		t.Fatal("agent2 Claim() succeeded, want error")
	}

	got, err := taskStore.Get("task-concurrent")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.AgentID != "agent-1" {
		t.Errorf("AgentID = %q, want %q", got.AgentID, "agent-1")
	}

	if got.Version != 1 {
		t.Errorf("Version = %d, want %d", got.Version, 1)
	}
}

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

func TestAgentClaimSetsLeaseAndFencingToken(t *testing.T) {
	taskStore := store.NewTaskStore()
	agent := NewAgent("agent-1", taskStore)

	now := time.Now()

	input := task.Task{
		ID:        "task-lease",
		Payload:   "process order",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := taskStore.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	beforeClaim := time.Now()

	if err := agent.Claim("task-lease"); err != nil {
		t.Fatalf("Claim() error = %v", err)
	}

	afterClaim := time.Now()

	got, err := taskStore.Get("task-lease")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.LeaseExpiresAt == nil {
		t.Fatal("LeaseExpiresAt is nil")
	}

	if got.LeaseExpiresAt.Before(beforeClaim.Add(LeaseDuration)) {
		t.Errorf("LeaseExpiresAt = %v, expected at least around %v",
			got.LeaseExpiresAt,
			beforeClaim.Add(LeaseDuration))
	}

	if got.LeaseExpiresAt.After(afterClaim.Add(LeaseDuration)) {
		t.Errorf("LeaseExpiresAt = %v, expected around %v",
			got.LeaseExpiresAt,
			afterClaim.Add(LeaseDuration))
	}

	if got.FencingToken != 1 {
		t.Errorf("FencingToken = %d, want %d", got.FencingToken, 1)
	}
}

func TestAgentRenew(t *testing.T) {
	taskStore := store.NewTaskStore()
	agent := NewAgent("agent-1", taskStore)

	now := time.Now()

	input := task.Task{
		ID:        "task-renew",
		Payload:   "process order",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := taskStore.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := agent.Claim("task-renew"); err != nil {
		t.Fatalf("Claim() error = %v", err)
	}

	beforeRenew, err := taskStore.Get("task-renew")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if beforeRenew.LeaseExpiresAt == nil {
		t.Fatal("LeaseExpiresAt is nil after Claim()")
	}

	oldExpiry := *beforeRenew.LeaseExpiresAt
	oldToken := beforeRenew.FencingToken
	oldVersion := beforeRenew.Version

	if err := agent.Renew("task-renew"); err != nil {
		t.Fatalf("Renew() error = %v", err)
	}

	afterRenew, err := taskStore.Get("task-renew")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if afterRenew.LeaseExpiresAt == nil {
		t.Fatal("LeaseExpiresAt is nil after Renew()")
	}

	if !afterRenew.LeaseExpiresAt.After(oldExpiry) {
		t.Errorf("LeaseExpiresAt did not move forward: old=%v new=%v",
			oldExpiry, *afterRenew.LeaseExpiresAt)
	}

	if afterRenew.FencingToken != oldToken {
		t.Errorf("FencingToken = %d, want %d",
			afterRenew.FencingToken, oldToken)
	}

	if afterRenew.Version != oldVersion+1 {
		t.Errorf("Version = %d, want %d",
			afterRenew.Version, oldVersion+1)
	}
}

func TestAgentRecoverExpiredTask(t *testing.T) {
	taskStore := store.NewTaskStore()

	agent1 := NewAgent("agent-1", taskStore)
	agent2 := NewAgent("agent-2", taskStore)

	now := time.Now()

	input := task.Task{
		ID:        "task-recovery",
		Payload:   "process order",
		Status:    task.StatusPending,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := taskStore.Create(input); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := agent1.Claim("task-recovery"); err != nil {
		t.Fatalf("agent1 Claim() error = %v", err)
	}

	claimed, err := taskStore.Get("task-recovery")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	// Simulate Agent 1's lease expiring.
	expired := time.Now().Add(-time.Second)
	claimed.LeaseExpiresAt = &expired

	if err := taskStore.Update(claimed, claimed.Version); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if err := agent2.Recover("task-recovery"); err != nil {
		t.Fatalf("agent2 Recover() error = %v", err)
	}

	recovered, err := taskStore.Get("task-recovery")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if recovered.AgentID != "agent-2" {
		t.Errorf("AgentID = %q, want %q", recovered.AgentID, "agent-2")
	}

	if recovered.Status != task.StatusClaimed {
		t.Errorf("Status = %q, want %q", recovered.Status, task.StatusClaimed)
	}

	if recovered.FencingToken != 2 {
		t.Errorf("FencingToken = %d, want %d", recovered.FencingToken, 2)
	}

	if recovered.Version != 3 {
		t.Errorf("Version = %d, want %d", recovered.Version, 3)
	}

	if recovered.LeaseExpiresAt == nil {
		t.Fatal("LeaseExpiresAt is nil after recovery")
	}

	if recovered.IsLeaseExpired(time.Now()) {
		t.Error("recovered task lease should be active")
	}
}

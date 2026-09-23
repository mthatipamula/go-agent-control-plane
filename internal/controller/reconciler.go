package controller

import (
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

type Reconciler struct {
	store *store.TaskStore
}

func NewReconciler(store *store.TaskStore) *Reconciler {
	return &Reconciler{
		store: store,
	}
}

func (r *Reconciler) Reconcile(taskID string) error {
	t, err := r.store.Get(taskID)
	if err != nil {
		return err
	}

	// Only claimed or running tasks can become stale.
	if t.Status != task.StatusClaimed && t.Status != task.StatusRunning {
		return nil
	}

	// The task is still healthy; nothing to do.
	if !t.IsLeaseExpired(time.Now()) {
		return nil
	}

	// The current lease expired, so make the task available
	// for another agent to claim.
	t.Status = task.StatusPending
	t.AgentID = ""
	t.LeaseExpiresAt = nil
	t.UpdatedAt = time.Now()

	return r.store.Update(t, t.Version)
}
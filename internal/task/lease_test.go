package task

import (
	"testing"
	"time"
)

func TestIsLeaseExpired(t *testing.T) {
	expiry := time.Now().Add(30 * time.Second)

	task := Task{
		ID:             "task-lease",
		LeaseExpiresAt: &expiry,
	}

	if task.IsLeaseExpired(time.Now()) {
		t.Fatal("lease should not be expired yet")
	}

	if !task.IsLeaseExpired(expiry) {
		t.Fatal("lease should be expired at the expiration time")
	}

	if !task.IsLeaseExpired(expiry.Add(time.Second)) {
		t.Fatal("lease should be expired after the expiration time")
	}
}

func TestIsLeaseExpiredWithoutLease(t *testing.T) {
	task := Task{
		ID: "task-no-lease",
	}

	if task.IsLeaseExpired(time.Now()) {
		t.Fatal("task without a lease should not be expired")
	}
}

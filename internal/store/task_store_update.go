package store

import (
	"errors"

	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

var ErrVersionConflict = errors.New("task version conflict")

func (s *TaskStore) Update(t task.Task, expectedVersion int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.tasks[t.ID]
	if !exists {
		return ErrTaskNotFound
	}

	if current.Version != expectedVersion {
		return ErrVersionConflict
	}

	t.Version = expectedVersion + 1
	s.tasks[t.ID] = t

	return nil
}

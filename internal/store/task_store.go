package store

import (
	"errors"
	"sync"

	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]task.Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]task.Task),
	}
}

func (s *TaskStore) Create(t task.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[t.ID]; exists {
		return errors.New("task already exists")
	}

	s.tasks[t.ID] = t
	return nil
}

func (s *TaskStore) Get(id string) (task.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, exists := s.tasks[id]
	if !exists {
		return task.Task{}, ErrTaskNotFound
	}

	return t, nil
}

func (s *TaskStore) UpdateWithFencing(
	t task.Task,
	expectedVersion int64,
	expectedFencingToken int64,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.tasks[t.ID]
	if !exists {
		return ErrTaskNotFound
	}

	if current.Version != expectedVersion {
		return ErrVersionConflict
	}

	if current.FencingToken != expectedFencingToken {
		return ErrFencingTokenConflict
	}

	t.Version = expectedVersion + 1
	s.tasks[t.ID] = t

	return nil
}

func (s *TaskStore) ListPending() []task.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var pending []task.Task

	for _, t := range s.tasks {
		if t.Status == task.StatusPending {
			pending = append(pending, t)
		}
	}

	return pending
}

func (s *TaskStore) List() []task.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]task.Task, 0, len(s.tasks))

	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}

	return tasks
}
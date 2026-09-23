package store

import (
	"database/sql"
	"fmt"

	"github.com/mthatipamula/go-agent-control-plane/internal/task"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(db *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{
		db: db,
	}
}

func OpenPostgres(dsn string) (*sql.DB, error) {
	return sql.Open("pgx", dsn)
}

func (s *PostgresTaskStore) Create(t task.Task) error {
	_, err := s.db.Exec(`
		INSERT INTO tasks (
			id,
			payload,
			status,
			agent_id,
			attempt,
			version,
			created_at,
			updated_at,
			lease_expires_at,
			fencing_token
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		t.ID,
		t.Payload,
		t.Status,
		t.AgentID,
		t.Attempt,
		t.Version,
		t.CreatedAt,
		t.UpdatedAt,
		t.LeaseExpiresAt,
		t.FencingToken,
	)

	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (s *PostgresTaskStore) Get(id string) (task.Task, error) {
	var t task.Task

	err := s.db.QueryRow(`
		SELECT
			id,
			payload,
			status,
			agent_id,
			attempt,
			version,
			created_at,
			updated_at,
			lease_expires_at,
			fencing_token
		FROM tasks
		WHERE id = $1
	`, id).Scan(
		&t.ID,
		&t.Payload,
		&t.Status,
		&t.AgentID,
		&t.Attempt,
		&t.Version,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.LeaseExpiresAt,
		&t.FencingToken,
	)

	if err != nil {
		return task.Task{}, fmt.Errorf("get task %s: %w", id, err)
	}

	return t, nil
}

func (s *PostgresTaskStore) UpdateWithFencing(
	t task.Task,
	expectedVersion int64,
	expectedFencingToken int64,
) error {
	result, err := s.db.Exec(`
		UPDATE tasks
		SET
			payload = $1,
			status = $2,
			agent_id = $3,
			attempt = $4,
			version = $5,
			created_at = $6,
			updated_at = $7,
			lease_expires_at = $8,
			fencing_token = $9
		WHERE id = $10
		  AND version = $11
		  AND fencing_token = $12
	`,
		t.Payload,
		t.Status,
		t.AgentID,
		t.Attempt,
		expectedVersion+1,
		t.CreatedAt,
		t.UpdatedAt,
		t.LeaseExpiresAt,
		t.FencingToken,
		t.ID,
		expectedVersion,
		expectedFencingToken,
	)
	if err != nil {
		return fmt.Errorf("update task %s: %w", t.ID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check task update %s: %w", t.ID, err)
	}

	if rows == 0 {
		return ErrVersionConflict
	}

	return nil
}

func (s *PostgresTaskStore) List() ([]task.Task, error) {
	rows, err := s.db.Query(`
		SELECT
			id,
			payload,
			status,
			agent_id,
			attempt,
			version,
			created_at,
			updated_at,
			lease_expires_at,
			fencing_token
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]task.Task, 0)

	for rows.Next() {
		var t task.Task

		if err := rows.Scan(
			&t.ID,
			&t.Payload,
			&t.Status,
			&t.AgentID,
			&t.Attempt,
			&t.Version,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.LeaseExpiresAt,
			&t.FencingToken,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (s *PostgresTaskStore) ListPending() ([]task.Task, error) {
	rows, err := s.db.Query(`
		SELECT
			id,
			payload,
			status,
			agent_id,
			attempt,
			version,
			created_at,
			updated_at,
			lease_expires_at,
			fencing_token
		FROM tasks
		WHERE status = $1
		ORDER BY created_at ASC
	`, task.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("list pending tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]task.Task, 0)

	for rows.Next() {
		var t task.Task

		if err := rows.Scan(
			&t.ID,
			&t.Payload,
			&t.Status,
			&t.AgentID,
			&t.Attempt,
			&t.Version,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.LeaseExpiresAt,
			&t.FencingToken,
		); err != nil {
			return nil, fmt.Errorf("scan pending task: %w", err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending tasks: %w", err)
	}

	return tasks, nil
}

package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/todoist/backend/internal/db"
	"github.com/todoist/backend/internal/domain"
)

type TaskRepo struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{pool: pool}
}

func (r *TaskRepo) Create(ctx context.Context, t *domain.Task) error {
	const q = `
		INSERT INTO tasks (
			id, user_id, project_id, parent_task_id,
			title, description_md, description_html, links,
			priority, due_at, completed_at, sort_order,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8::jsonb,
			$9, $10, $11, $12,
			$13, $14
		)`

	linksJSON, err := json.Marshal(t.Links)
	if err != nil {
		return fmt.Errorf("marshal links: %w", err)
	}

	_, err = db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		t.ID, t.UserID, t.ProjectID, t.ParentTaskID,
		t.Title, t.DescriptionMD, t.DescriptionHTML, string(linksJSON),
		t.Priority, t.DueAt, t.CompletedAt, t.SortOrder,
		t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (r *TaskRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Task, error) {
	const q = `
		SELECT
			id, user_id, project_id, parent_task_id,
			title, description_md, description_html, links,
			priority, due_at, completed_at, sort_order,
			created_at, updated_at
		FROM tasks
		WHERE id = $1 AND user_id = $2`

	t, err := scanTask(db.ExtractDB(ctx, r.pool).QueryRow(ctx, q, id, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get task: %w", err)
	}
	return t, nil
}

func (r *TaskRepo) List(ctx context.Context, userID uuid.UUID, f domain.TaskFilter) ([]*domain.Task, error) {
	query, args := buildListQuery(userID, f)

	rows, err := db.ExtractDB(ctx, r.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepo) Update(ctx context.Context, t *domain.Task) error {
	const q = `
		UPDATE tasks
		SET
			project_id = $1, parent_task_id = $2,
			title = $3, description_md = $4, description_html = $5, links = $6::jsonb,
			priority = $7, due_at = $8, completed_at = $9, sort_order = $10,
			updated_at = $11
		WHERE id = $12 AND user_id = $13`

	linksJSON, err := json.Marshal(t.Links)
	if err != nil {
		return fmt.Errorf("marshal links: %w", err)
	}

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q,
		t.ProjectID, t.ParentTaskID,
		t.Title, t.DescriptionMD, t.DescriptionHTML, string(linksJSON),
		t.Priority, t.DueAt, t.CompletedAt, t.SortOrder,
		time.Now(), t.ID, t.UserID,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	const q = `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	tag, err := db.ExtractDB(ctx, r.pool).Exec(ctx, q, id, userID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// SetLabels replaces all labels on a task atomically.
func (r *TaskRepo) SetLabels(ctx context.Context, taskID, userID uuid.UUID, labelIDs []uuid.UUID) error {
	dbtx := db.ExtractDB(ctx, r.pool)

	if _, err := dbtx.Exec(ctx, `DELETE FROM task_labels WHERE task_id = $1`, taskID); err != nil {
		return fmt.Errorf("clear task labels: %w", err)
	}
	if len(labelIDs) == 0 {
		return nil
	}

	// Build a multi-row INSERT in one round-trip.
	placeholders := make([]string, len(labelIDs))
	args := make([]any, 0, len(labelIDs)*3)
	for i, id := range labelIDs {
		base := i*3 + 1
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d)", base, base+1, base+2)
		args = append(args, taskID, id, userID)
	}
	q := `INSERT INTO task_labels (task_id, label_id, user_id) VALUES ` + strings.Join(placeholders, ", ")

	if _, err := dbtx.Exec(ctx, q, args...); err != nil {
		return fmt.Errorf("insert task labels: %w", err)
	}
	return nil
}

// GetLabels returns all labels attached to a task.
func (r *TaskRepo) GetLabels(ctx context.Context, taskID uuid.UUID) ([]*domain.Label, error) {
	const q = `
		SELECT l.id, l.user_id, l.name, l.color, l.created_at, l.updated_at
		FROM labels l
		JOIN task_labels tl ON tl.label_id = l.id
		WHERE tl.task_id = $1
		ORDER BY l.name ASC`

	rows, err := db.ExtractDB(ctx, r.pool).Query(ctx, q, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task labels: %w", err)
	}
	defer rows.Close()

	var labels []*domain.Label
	for rows.Next() {
		var l domain.Label
		if err := rows.Scan(&l.ID, &l.UserID, &l.Name, &l.Color, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan label: %w", err)
		}
		labels = append(labels, &l)
	}
	return labels, rows.Err()
}

// buildListQuery constructs a parameterized SELECT for the task list endpoint.
// All user-controlled values go through numbered parameters ($N) — no untrusted
// input is ever interpolated directly into the query string.
func buildListQuery(userID uuid.UUID, f domain.TaskFilter) (string, []any) {
	args := []any{userID}
	conds := []string{"t.user_id = $1"}

	// addCond appends val to args and formats cond with the resulting parameter index.
	addCond := func(cond string, val any) {
		args = append(args, val)
		conds = append(conds, fmt.Sprintf(cond, len(args)))
	}

	if f.ProjectID != nil {
		addCond("t.project_id = $%d", *f.ProjectID)
	}
	if f.Priority != nil {
		addCond("t.priority = $%d", *f.Priority)
	}
	if f.DueBefore != nil {
		addCond("t.due_at <= $%d", *f.DueBefore)
	}
	if f.DueAfter != nil {
		addCond("t.due_at >= $%d", *f.DueAfter)
	}
	if f.ParentID != nil {
		addCond("t.parent_task_id = $%d", *f.ParentID)
	}
	if f.Completed != nil {
		if *f.Completed {
			conds = append(conds, "t.completed_at IS NOT NULL")
		} else {
			conds = append(conds, "t.completed_at IS NULL")
		}
	}
	if f.Search != "" {
		addCond("t.title ILIKE '%%' || $%d || '%%'", f.Search)
	}

	join := ""
	if f.LabelID != nil {
		join = "JOIN task_labels tl ON tl.task_id = t.id"
		addCond("tl.label_id = $%d", *f.LabelID)
	}

	limit := f.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	offset := max(f.Offset, 0)

	args = append(args, limit, offset)
	limitIdx, offsetIdx := len(args)-1, len(args)

	q := fmt.Sprintf(`
		SELECT DISTINCT
			t.id, t.user_id, t.project_id, t.parent_task_id,
			t.title, t.description_md, t.description_html, t.links,
			t.priority, t.due_at, t.completed_at, t.sort_order,
			t.created_at, t.updated_at
		FROM tasks t
		%s
		WHERE %s
		ORDER BY t.sort_order ASC, t.created_at ASC
		LIMIT $%d OFFSET $%d`,
		join,
		strings.Join(conds, " AND "),
		limitIdx,
		offsetIdx,
	)

	return q, args
}

// scanTask reads a task row from anything that has a Scan method.
func scanTask(row interface{ Scan(...any) error }) (*domain.Task, error) {
	var t domain.Task
	var linksRaw []byte

	err := row.Scan(
		&t.ID, &t.UserID, &t.ProjectID, &t.ParentTaskID,
		&t.Title, &t.DescriptionMD, &t.DescriptionHTML, &linksRaw,
		&t.Priority, &t.DueAt, &t.CompletedAt, &t.SortOrder,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(linksRaw) > 0 {
		if err := json.Unmarshal(linksRaw, &t.Links); err != nil {
			return nil, fmt.Errorf("unmarshal links: %w", err)
		}
	}
	return &t, nil
}

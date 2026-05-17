package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type AssignedTaskRepository struct {
	db *sql.DB
}

func NewAssignedTaskRepository() *AssignedTaskRepository {
	return &AssignedTaskRepository{db: database.DB}
}

func (r *AssignedTaskRepository) Create(task *models.AssignedTask) error {
	query := `
		INSERT INTO assigned_tasks (
			id, instance_id, step_id, step_name, assigned_to, assigned_by,
			status, due_date, org_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at`

	task.ID = uuid.New().String()
	var orgID interface{}
	if task.OrgID != "" {
		orgID = task.OrgID
	}
	return r.db.QueryRow(query,
		task.ID, task.InstanceID, task.StepID, task.StepName,
		task.AssignedTo, task.AssignedBy, task.Status, task.DueDate, orgID,
	).Scan(&task.CreatedAt, &task.UpdatedAt)
}

func (r *AssignedTaskRepository) GetActiveTaskForInstance(instanceID, orgID string) (*models.AssignedTask, error) {
	query := `
		SELECT at.id, at.instance_id, at.step_id, at.step_name,
		       at.assigned_to, at.assigned_by, at.status,
		       at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details
		FROM assigned_tasks at
		JOIN workflow_instances wi ON at.instance_id = wi.id
		WHERE at.instance_id = $1 AND at.org_id = $2
		  AND at.status IN ('pending', 'in_progress')
		ORDER BY at.created_at DESC
		LIMIT 1`

	return r.scanTask(r.db.QueryRow(query, instanceID, orgID))
}

func (r *AssignedTaskRepository) Complete(taskID, orgID string) error {
	query := `
		UPDATE assigned_tasks
		SET status = 'completed', completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND org_id = $2`
	result, err := r.db.Exec(query, taskID, orgID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *AssignedTaskRepository) GetByAssignee(orgID, assigneeID string, statusFilter ...string) ([]models.AssignedTask, error) {
	args := []interface{}{assigneeID, orgID}
	where := "at.assigned_to = $1 AND at.org_id = $2"
	i := 3

	if len(statusFilter) > 0 && statusFilter[0] != "" {
		where += fmt.Sprintf(" AND at.status = $%d", i)
		args = append(args, statusFilter[0])
	}

	query := fmt.Sprintf(`
		SELECT at.id, at.instance_id, at.step_id, at.step_name,
		       at.assigned_to, at.assigned_by, at.status,
		       at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details
		FROM assigned_tasks at
		JOIN workflow_instances wi ON at.instance_id = wi.id
		WHERE %s
		ORDER BY at.created_at DESC`, where)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.AssignedTask
	for rows.Next() {
		t, err := r.scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

type taskRowScanner interface {
	Scan(dest ...interface{}) error
}

func (r *AssignedTaskRepository) scanTask(row taskRowScanner) (*models.AssignedTask, error) {
	var t models.AssignedTask
	var taskDetailsJSON []byte
	var dueDate, completedAt sql.NullTime

	err := row.Scan(
		&t.ID, &t.InstanceID, &t.StepID, &t.StepName,
		&t.AssignedTo, &t.AssignedBy, &t.Status,
		&dueDate, &completedAt, &t.CreatedAt, &t.UpdatedAt,
		&taskDetailsJSON,
	)
	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}
	if completedAt.Valid {
		t.CompletedAt = &completedAt.Time
	}

	if len(taskDetailsJSON) > 0 {
		var td models.TaskDetails
		if err := json.Unmarshal(taskDetailsJSON, &td); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task_details: %w", err)
		}
		t.TaskDetails = &td
	}

	return &t, nil
}

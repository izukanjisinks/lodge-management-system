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
			status, due_date, org_id, branch_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`

	task.ID = uuid.New().String()
	var orgID interface{}
	if task.OrgID != "" {
		orgID = task.OrgID
	}
	return r.db.QueryRow(query,
		task.ID, task.InstanceID, task.StepID, task.StepName,
		task.AssignedTo, task.AssignedBy, task.Status, task.DueDate, orgID, task.BranchID,
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

// GetByID returns a single task by ID, scoped to the org (any assignee). Used by
// the task-detail view, where managers/admins may open tasks assigned to others.
func (r *AssignedTaskRepository) GetByID(taskID, orgID string) (*models.AssignedTask, error) {
	query := `
		SELECT at.id, at.instance_id, at.step_id, at.step_name,
		       at.assigned_to, at.assigned_by, at.status,
		       at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details
		FROM assigned_tasks at
		JOIN workflow_instances wi ON at.instance_id = wi.id
		WHERE at.id = $1 AND at.org_id = $2`
	return r.scanTask(r.db.QueryRow(query, taskID, orgID))
}

// statusClause appends a status filter to the WHERE string. It accepts two group
// keywords — "active" (pending + in_progress) and "completed" (completed + rejected)
// — as well as any exact status value. Returns the updated clause, args, and next
// placeholder index. An empty statusFilter is a no-op.
func statusClause(where string, args []interface{}, i int, statusFilter string) (string, []interface{}, int) {
	switch statusFilter {
	case "":
		return where, args, i
	case "active":
		where += fmt.Sprintf(" AND at.status IN ($%d, $%d)", i, i+1)
		args = append(args, "pending", "in_progress")
		return where, args, i + 2
	case "completed":
		where += fmt.Sprintf(" AND at.status IN ($%d, $%d)", i, i+1)
		args = append(args, "completed", "rejected")
		return where, args, i + 2
	default:
		where += fmt.Sprintf(" AND at.status = $%d", i)
		args = append(args, statusFilter)
		return where, args, i + 1
	}
}

// GetByAssignee returns the tasks assigned to a user, paginated. It returns the
// page of tasks plus the total count matching the filter (for pagination UIs).
func (r *AssignedTaskRepository) GetByAssignee(orgID, assigneeID, statusFilter string, limit, offset int) ([]models.AssignedTask, int, error) {
	args := []interface{}{assigneeID, orgID}
	where := "at.assigned_to = $1 AND at.org_id = $2"
	i := 3

	where, args, i = statusClause(where, args, i, statusFilter)

	var total int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM assigned_tasks at
		JOIN workflow_instances wi ON at.instance_id = wi.id AND wi.status != 'cancelled'
		WHERE %s`, where)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT at.id, at.instance_id, at.step_id, at.step_name,
		       at.assigned_to, at.assigned_by, at.status,
		       at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details
		FROM assigned_tasks at
		JOIN workflow_instances wi ON at.instance_id = wi.id AND wi.status != 'cancelled'
		WHERE %s
		ORDER BY at.created_at DESC
		LIMIT $%d OFFSET $%d`, where, i, i+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []models.AssignedTask
	for rows.Next() {
		t, err := r.scanTask(rows)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, total, rows.Err()
}

// GetAllByOrg returns all tasks in an org regardless of assignee, paginated and
// with the assignee's display name joined from the users table. Returns the page
// of tasks plus the total count matching the filter.
func (r *AssignedTaskRepository) GetAllByOrg(orgID, branchID, statusFilter string, limit, offset int) ([]models.AssignedTask, int, error) {
	args := []interface{}{orgID}
	where := "at.org_id = $1"
	i := 2

	if branchID != "" {
		where += fmt.Sprintf(" AND at.branch_id = $%d", i)
		args = append(args, branchID)
		i++
	}
	where, args, i = statusClause(where, args, i, statusFilter)

	var total int
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM assigned_tasks at
		JOIN workflow_instances wi ON at.instance_id = wi.id AND wi.status != 'cancelled'
		WHERE %s`, where)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT at.id, at.instance_id, at.step_id, at.step_name,
		       at.assigned_to, at.assigned_by, at.status,
		       at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details,
		       COALESCE(u.full_name, u.email, at.assigned_to::text) AS assignee_name
		FROM assigned_tasks at
		JOIN workflow_instances wi ON at.instance_id = wi.id AND wi.status != 'cancelled'
		LEFT JOIN users u ON u.user_id = at.assigned_to
		WHERE %s
		ORDER BY at.created_at DESC
		LIMIT $%d OFFSET $%d`, where, i, i+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []models.AssignedTask
	for rows.Next() {
		var t models.AssignedTask
		var taskDetailsJSON []byte
		var dueDate, completedAt sql.NullTime
		var assigneeName sql.NullString

		err := rows.Scan(
			&t.ID, &t.InstanceID, &t.StepID, &t.StepName,
			&t.AssignedTo, &t.AssignedBy, &t.Status,
			&dueDate, &completedAt, &t.CreatedAt, &t.UpdatedAt,
			&taskDetailsJSON,
			&assigneeName,
		)
		if err != nil {
			return nil, 0, err
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		if completedAt.Valid {
			t.CompletedAt = &completedAt.Time
		}
		if assigneeName.Valid {
			t.AssigneeName = assigneeName.String
		}
		if len(taskDetailsJSON) > 0 {
			var td models.TaskDetails
			if err := json.Unmarshal(taskDetailsJSON, &td); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal task_details: %w", err)
			}
			t.TaskDetails = &td
		}
		tasks = append(tasks, t)
	}
	return tasks, total, rows.Err()
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

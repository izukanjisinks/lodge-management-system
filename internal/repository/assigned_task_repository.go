package repository

import (
	"database/sql"
	"encoding/json"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type AssignedTaskRepository struct {
	db *sql.DB
}

func NewAssignedTaskRepository() *AssignedTaskRepository {
	return &AssignedTaskRepository{
		db: database.DB,
	}
}

// Create creates a new assigned task
func (r *AssignedTaskRepository) Create(task *models.AssignedTask) error {
	query := `
		INSERT INTO assigned_tasks (
			id, instance_id, step_id, step_name, assigned_to,
			assigned_by, status, due_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	task.ID = uuid.New().String()

	return r.db.QueryRow(
		query,
		task.ID,
		task.InstanceID,
		task.StepID,
		task.StepName,
		task.AssignedTo,
		task.AssignedBy,
		task.Status,
		task.DueDate,
	).Scan(&task.CreatedAt, &task.UpdatedAt)
}

// GetByID retrieves a task by ID
func (r *AssignedTaskRepository) GetByID(id string) (*models.AssignedTask, error) {
	query := `
		SELECT at.id, at.instance_id, at.step_id, at.step_name, at.assigned_to, at.assigned_by,
		       at.status, at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details
		FROM assigned_tasks at
		LEFT JOIN workflow_instances wi ON at.instance_id = wi.id
		WHERE at.id = $1
	`

	var task models.AssignedTask
	var taskDetailsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.InstanceID,
		&task.StepID,
		&task.StepName,
		&task.AssignedTo,
		&task.AssignedBy,
		&task.Status,
		&task.DueDate,
		&task.CompletedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
		&taskDetailsJSON,
	)

	if err != nil {
		return nil, err
	}

	// Parse task_details JSON if present
	if len(taskDetailsJSON) > 0 {
		var details models.TaskDetails
		if err := json.Unmarshal(taskDetailsJSON, &details); err == nil {
			task.TaskDetails = &details
		}
	}

	return &task, nil
}

// GetByAssignee retrieves all tasks assigned to a user
func (r *AssignedTaskRepository) GetByAssignee(userID string, status ...string) ([]models.AssignedTask, error) {
	var query string
	var args []interface{}

	if len(status) > 0 {
		query = `
			SELECT at.id, at.instance_id, at.step_id, at.step_name, at.assigned_to, at.assigned_by,
			       at.status, at.due_date, at.completed_at, at.created_at, at.updated_at,
			       wi.task_details
			FROM assigned_tasks at
			LEFT JOIN workflow_instances wi ON at.instance_id = wi.id
			WHERE at.assigned_to = $1 AND at.status = $2
			ORDER BY at.created_at DESC
		`
		args = []interface{}{userID, status[0]}
	} else {
		query = `
			SELECT at.id, at.instance_id, at.step_id, at.step_name, at.assigned_to, at.assigned_by,
			       at.status, at.due_date, at.completed_at, at.created_at, at.updated_at,
			       wi.task_details
			FROM assigned_tasks at
			LEFT JOIN workflow_instances wi ON at.instance_id = wi.id
			WHERE at.assigned_to = $1
			ORDER BY at.created_at DESC
		`
		args = []interface{}{userID}
	}

	return r.queryTasks(query, args...)
}

// GetPendingByAssignee retrieves all pending tasks for a user
func (r *AssignedTaskRepository) GetPendingByAssignee(userID string) ([]models.AssignedTask, error) {
	return r.GetByAssignee(userID, "pending")
}

// GetByInstance retrieves all tasks for a workflow instance
func (r *AssignedTaskRepository) GetByInstance(instanceID string) ([]models.AssignedTask, error) {
	query := `
		SELECT at.id, at.instance_id, at.step_id, at.step_name, at.assigned_to, at.assigned_by,
		       at.status, at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details
		FROM assigned_tasks at
		LEFT JOIN workflow_instances wi ON at.instance_id = wi.id
		WHERE at.instance_id = $1
		ORDER BY at.created_at
	`

	return r.queryTasks(query, instanceID)
}

// GetActiveTaskForInstance retrieves the current active task for an instance
func (r *AssignedTaskRepository) GetActiveTaskForInstance(instanceID string) (*models.AssignedTask, error) {
	query := `
		SELECT at.id, at.instance_id, at.step_id, at.step_name, at.assigned_to, at.assigned_by,
		       at.status, at.due_date, at.completed_at, at.created_at, at.updated_at,
		       wi.task_details
		FROM assigned_tasks at
		LEFT JOIN workflow_instances wi ON at.instance_id = wi.id
		WHERE at.instance_id = $1 AND at.status IN ('pending', 'in_progress')
		ORDER BY at.created_at DESC
		LIMIT 1
	`

	var task models.AssignedTask
	var taskDetailsJSON []byte

	err := r.db.QueryRow(query, instanceID).Scan(
		&task.ID,
		&task.InstanceID,
		&task.StepID,
		&task.StepName,
		&task.AssignedTo,
		&task.AssignedBy,
		&task.Status,
		&task.DueDate,
		&task.CompletedAt,
		&task.CreatedAt,
		&task.UpdatedAt,
		&taskDetailsJSON,
	)

	if err != nil {
		return nil, err
	}

	// Parse task_details JSON if present
	if len(taskDetailsJSON) > 0 {
		var details models.TaskDetails
		if err := json.Unmarshal(taskDetailsJSON, &details); err == nil {
			task.TaskDetails = &details
		}
	}

	return &task, nil
}

// UpdateStatus updates the status of a task
func (r *AssignedTaskRepository) UpdateStatus(taskID, status string) error {
	query := `
		UPDATE assigned_tasks
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(query, status, taskID)
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

// Complete marks a task as completed
func (r *AssignedTaskRepository) Complete(taskID string) error {
	query := `
		UPDATE assigned_tasks
		SET status = 'completed',
		    completed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query, taskID)
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

// Reassign reassigns a task to a different user
func (r *AssignedTaskRepository) Reassign(taskID, newAssigneeID, reassignedByID string) error {
	query := `
		UPDATE assigned_tasks
		SET assigned_to = $1,
		    assigned_by = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.Exec(query, newAssigneeID, reassignedByID, taskID)
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

// CountPendingByAssignee counts pending tasks for a user
func (r *AssignedTaskRepository) CountPendingByAssignee(userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM assigned_tasks
		WHERE assigned_to = $1 AND status = 'pending'
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Helper method to query multiple tasks
func (r *AssignedTaskRepository) queryTasks(query string, args ...interface{}) ([]models.AssignedTask, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.AssignedTask
	for rows.Next() {
		var task models.AssignedTask
		var taskDetailsJSON []byte

		err := rows.Scan(
			&task.ID,
			&task.InstanceID,
			&task.StepID,
			&task.StepName,
			&task.AssignedTo,
			&task.AssignedBy,
			&task.Status,
			&task.DueDate,
			&task.CompletedAt,
			&task.CreatedAt,
			&task.UpdatedAt,
			&taskDetailsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Parse task_details JSON if present
		if len(taskDetailsJSON) > 0 {
			var details models.TaskDetails
			if err := json.Unmarshal(taskDetailsJSON, &details); err == nil {
				task.TaskDetails = &details
			}
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
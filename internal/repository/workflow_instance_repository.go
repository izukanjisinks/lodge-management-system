package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type WorkflowInstanceRepository struct {
	db *sql.DB
}

func NewWorkflowInstanceRepository() *WorkflowInstanceRepository {
	return &WorkflowInstanceRepository{
		db: database.DB,
	}
}

// Create creates a new workflow instance
func (r *WorkflowInstanceRepository) Create(instance *models.WorkflowInstance) error {
	query := `
		INSERT INTO workflow_instances (
			id, workflow_id, current_step_id, status, task_details,
			created_by, due_date, priority
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	instance.ID = uuid.New().String()

	// Marshal task_details to JSON
	taskDetailsJSON, err := json.Marshal(instance.TaskDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal task_details: %w", err)
	}

	return r.db.QueryRow(
		query,
		instance.ID,
		instance.WorkflowID,
		instance.CurrentStepID,
		instance.Status,
		taskDetailsJSON,
		instance.CreatedBy,
		instance.DueDate,
		instance.Priority,
	).Scan(&instance.CreatedAt, &instance.UpdatedAt)
}

// GetByID retrieves a workflow instance by ID
func (r *WorkflowInstanceRepository) GetByID(id string) (*models.WorkflowInstance, error) {
	query := `
		SELECT id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE id = $1
	`

	var instance models.WorkflowInstance
	var taskDetailsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&instance.ID,
		&instance.WorkflowID,
		&instance.CurrentStepID,
		&instance.Status,
		&taskDetailsJSON,
		&instance.CreatedBy,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.CompletedAt,
		&instance.DueDate,
		&instance.Priority,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal task_details
	if err := json.Unmarshal(taskDetailsJSON, &instance.TaskDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task_details: %w", err)
	}

	return &instance, nil
}

// GetByTaskID retrieves a workflow instance by task ID (from task_details)
func (r *WorkflowInstanceRepository) GetByTaskID(taskID string) (*models.WorkflowInstance, error) {
	query := `
		SELECT id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE task_details->>'task_id' = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var instance models.WorkflowInstance
	var taskDetailsJSON []byte

	err := r.db.QueryRow(query, taskID).Scan(
		&instance.ID,
		&instance.WorkflowID,
		&instance.CurrentStepID,
		&instance.Status,
		&taskDetailsJSON,
		&instance.CreatedBy,
		&instance.CreatedAt,
		&instance.UpdatedAt,
		&instance.CompletedAt,
		&instance.DueDate,
		&instance.Priority,
	)

	if err != nil {
		return nil, err
	}

	// Unmarshal task_details
	if err := json.Unmarshal(taskDetailsJSON, &instance.TaskDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task_details: %w", err)
	}

	return &instance, nil
}

// GetByCreator retrieves all workflow instances created by a user
func (r *WorkflowInstanceRepository) GetByCreator(creatorID string) ([]models.WorkflowInstance, error) {
	query := `
		SELECT id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE created_by = $1
		ORDER BY created_at DESC
	`

	return r.queryInstances(query, creatorID)
}

// GetByStatus retrieves workflow instances by status
func (r *WorkflowInstanceRepository) GetByStatus(status string) ([]models.WorkflowInstance, error) {
	query := `
		SELECT id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE status = $1
		ORDER BY created_at DESC
	`

	return r.queryInstances(query, status)
}

// UpdateStep updates the current step of a workflow instance
func (r *WorkflowInstanceRepository) UpdateStep(instanceID, newStepID, status string) error {
	query := `
		UPDATE workflow_instances
		SET current_step_id = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.Exec(query, newStepID, status, instanceID)
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

// Complete marks a workflow instance as completed
func (r *WorkflowInstanceRepository) Complete(instanceID string) error {
	query := `
		UPDATE workflow_instances
		SET status = 'completed',
		    completed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query, instanceID)
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

// Cancel marks a workflow instance as cancelled
func (r *WorkflowInstanceRepository) Cancel(instanceID string) error {
	query := `
		UPDATE workflow_instances
		SET status = 'cancelled',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query, instanceID)
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

// Helper method to query multiple instances
func (r *WorkflowInstanceRepository) queryInstances(query string, args ...interface{}) ([]models.WorkflowInstance, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []models.WorkflowInstance
	for rows.Next() {
		var instance models.WorkflowInstance
		var taskDetailsJSON []byte

		err := rows.Scan(
			&instance.ID,
			&instance.WorkflowID,
			&instance.CurrentStepID,
			&instance.Status,
			&taskDetailsJSON,
			&instance.CreatedBy,
			&instance.CreatedAt,
			&instance.UpdatedAt,
			&instance.CompletedAt,
			&instance.DueDate,
			&instance.Priority,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal task_details
		if err := json.Unmarshal(taskDetailsJSON, &instance.TaskDetails); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task_details: %w", err)
		}

		instances = append(instances, instance)
	}

	return instances, nil
}
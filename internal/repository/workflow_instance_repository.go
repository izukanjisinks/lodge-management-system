package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

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
			created_by, due_date, priority, org_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	instance.ID = uuid.New().String()

	// Marshal task_details to JSON
	taskDetailsJSON, err := json.Marshal(instance.TaskDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal task_details: %w", err)
	}

	var orgID interface{}
	if instance.OrgID != "" {
		orgID = instance.OrgID
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
		orgID,
	).Scan(&instance.CreatedAt, &instance.UpdatedAt)
}

// GetByID retrieves a workflow instance by ID, scoped to org.
func (r *WorkflowInstanceRepository) GetByID(id, orgID string) (*models.WorkflowInstance, error) {
	query := `
		SELECT id, org_id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE id = $1 AND org_id = $2
	`

	var instance models.WorkflowInstance
	var taskDetailsJSON []byte

	err := r.db.QueryRow(query, id, orgID).Scan(
		&instance.ID,
		&instance.OrgID,
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

// GetByTaskID retrieves a workflow instance by task ID (from task_details), scoped to org.
func (r *WorkflowInstanceRepository) GetByTaskID(taskID, orgID string) (*models.WorkflowInstance, error) {
	query := `
		SELECT id, org_id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE task_details->>'task_id' = $1 AND org_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var instance models.WorkflowInstance
	var taskDetailsJSON []byte

	err := r.db.QueryRow(query, taskID, orgID).Scan(
		&instance.ID,
		&instance.OrgID,
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

// GetByCreator retrieves all workflow instances created by a user, scoped to org.
func (r *WorkflowInstanceRepository) GetByCreator(orgID, creatorID string) ([]models.WorkflowInstance, error) {
	query := `
		SELECT id, org_id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE created_by = $1 AND org_id = $2
		ORDER BY created_at DESC
	`
	return r.queryInstances(query, creatorID, orgID)
}

// GetByStatus retrieves workflow instances by status, scoped to org.
func (r *WorkflowInstanceRepository) GetByStatus(orgID, status string) ([]models.WorkflowInstance, error) {
	query := `
		SELECT id, org_id, workflow_id, current_step_id, status, task_details,
		       created_by, created_at, updated_at, completed_at, due_date, priority
		FROM workflow_instances
		WHERE status = $1 AND org_id = $2
		ORDER BY created_at DESC
	`
	return r.queryInstances(query, status, orgID)
}

// UpdateStep updates the current step of a workflow instance, scoped to org.
func (r *WorkflowInstanceRepository) UpdateStep(instanceID, newStepID, status, orgID string) error {
	query := `
		UPDATE workflow_instances
		SET current_step_id = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND org_id = $4
	`

	result, err := r.db.Exec(query, newStepID, status, instanceID, orgID)
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

// Complete marks a workflow instance as completed, scoped to org.
func (r *WorkflowInstanceRepository) Complete(instanceID, orgID string) error {
	query := `
		UPDATE workflow_instances
		SET status = 'completed',
		    completed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND org_id = $2
	`

	result, err := r.db.Exec(query, instanceID, orgID)
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

// Cancel marks a workflow instance as cancelled, scoped to org.
func (r *WorkflowInstanceRepository) Cancel(instanceID, orgID string) error {
	query := `
		UPDATE workflow_instances
		SET status = 'cancelled',
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND org_id = $2
	`

	result, err := r.db.Exec(query, instanceID, orgID)
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
			&instance.OrgID,
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
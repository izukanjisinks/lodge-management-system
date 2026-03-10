package repository

import (
	"database/sql"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type WorkflowHistoryRepository struct {
	db *sql.DB
}

func NewWorkflowHistoryRepository() *WorkflowHistoryRepository {
	return &WorkflowHistoryRepository{
		db: database.DB,
	}
}

// Create creates a new workflow history entry
func (r *WorkflowHistoryRepository) Create(history *models.WorkflowHistory) error {
	query := `
		INSERT INTO workflow_history (
			id, instance_id, from_step_id, to_step_id, action_taken,
			performed_by, performed_by_name, comments, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING timestamp
	`

	history.ID = uuid.New().String()

	// Handle metadata - if empty, insert NULL instead of empty string (which is invalid JSON)
	var metadata interface{}
	if history.Metadata != "" {
		metadata = history.Metadata
	} else {
		metadata = nil
	}

	return r.db.QueryRow(
		query,
		history.ID,
		history.InstanceID,
		history.FromStepID,
		history.ToStepID,
		history.ActionTaken,
		history.PerformedBy,
		history.PerformedByName,
		history.Comments,
		metadata,
	).Scan(&history.Timestamp)
}

// GetByInstanceID retrieves all history entries for a workflow instance
func (r *WorkflowHistoryRepository) GetByInstanceID(instanceID string) ([]models.WorkflowHistory, error) {
	query := `
		SELECT id, instance_id, from_step_id, to_step_id, action_taken,
		       performed_by, performed_by_name, comments, metadata, timestamp
		FROM workflow_history
		WHERE instance_id = $1
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(query, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historyEntries []models.WorkflowHistory
	for rows.Next() {
		var entry models.WorkflowHistory
		err := rows.Scan(
			&entry.ID,
			&entry.InstanceID,
			&entry.FromStepID,
			&entry.ToStepID,
			&entry.ActionTaken,
			&entry.PerformedBy,
			&entry.PerformedByName,
			&entry.Comments,
			&entry.Metadata,
			&entry.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		historyEntries = append(historyEntries, entry)
	}

	return historyEntries, nil
}

// GetByPerformer retrieves all history entries performed by a user
func (r *WorkflowHistoryRepository) GetByPerformer(userID string) ([]models.WorkflowHistory, error) {
	query := `
		SELECT id, instance_id, from_step_id, to_step_id, action_taken,
		       performed_by, performed_by_name, comments, metadata, timestamp
		FROM workflow_history
		WHERE performed_by = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historyEntries []models.WorkflowHistory
	for rows.Next() {
		var entry models.WorkflowHistory
		err := rows.Scan(
			&entry.ID,
			&entry.InstanceID,
			&entry.FromStepID,
			&entry.ToStepID,
			&entry.ActionTaken,
			&entry.PerformedBy,
			&entry.PerformedByName,
			&entry.Comments,
			&entry.Metadata,
			&entry.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		historyEntries = append(historyEntries, entry)
	}

	return historyEntries, nil
}

// GetByAction retrieves all history entries for a specific action type
func (r *WorkflowHistoryRepository) GetByAction(action string) ([]models.WorkflowHistory, error) {
	query := `
		SELECT id, instance_id, from_step_id, to_step_id, action_taken,
		       performed_by, performed_by_name, comments, metadata, timestamp
		FROM workflow_history
		WHERE action_taken = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(query, action)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historyEntries []models.WorkflowHistory
	for rows.Next() {
		var entry models.WorkflowHistory
		err := rows.Scan(
			&entry.ID,
			&entry.InstanceID,
			&entry.FromStepID,
			&entry.ToStepID,
			&entry.ActionTaken,
			&entry.PerformedBy,
			&entry.PerformedByName,
			&entry.Comments,
			&entry.Metadata,
			&entry.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		historyEntries = append(historyEntries, entry)
	}

	return historyEntries, nil
}
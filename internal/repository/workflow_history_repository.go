package repository

import (
	"database/sql"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

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

// GetByPerformer retrieves all history entries performed by a user, scoped to org.
func (r *WorkflowHistoryRepository) GetByPerformer(orgID, userID string) ([]models.WorkflowHistory, error) {
	query := `
		SELECT wh.id, wh.instance_id, wh.from_step_id, wh.to_step_id, wh.action_taken,
		       wh.performed_by, wh.performed_by_name, wh.comments, wh.metadata, wh.timestamp
		FROM workflow_history wh
		JOIN workflow_instances wi ON wh.instance_id = wi.id
		WHERE wh.performed_by = $1 AND wi.org_id = $2
		ORDER BY wh.timestamp DESC
	`

	rows, err := r.db.Query(query, userID, orgID)
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

// GetByAction retrieves all history entries for a specific action type, scoped to org.
func (r *WorkflowHistoryRepository) GetByAction(orgID, action string) ([]models.WorkflowHistory, error) {
	query := `
		SELECT wh.id, wh.instance_id, wh.from_step_id, wh.to_step_id, wh.action_taken,
		       wh.performed_by, wh.performed_by_name, wh.comments, wh.metadata, wh.timestamp
		FROM workflow_history wh
		JOIN workflow_instances wi ON wh.instance_id = wi.id
		WHERE wh.action_taken = $1 AND wi.org_id = $2
		ORDER BY wh.timestamp DESC
	`

	rows, err := r.db.Query(query, action, orgID)
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
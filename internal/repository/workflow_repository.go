package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type WorkflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository() *WorkflowRepository {
	return &WorkflowRepository{
		db: database.DB,
	}
}

// GetByID retrieves a workflow template by ID
func (r *WorkflowRepository) GetByID(id string) (*models.Workflow, error) {
	query := `
		SELECT id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows
		WHERE id = $1
	`

	var wf models.Workflow
	var workflowType sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&wf.ID,
		&wf.Name,
		&wf.Description,
		&workflowType,
		&wf.IsActive,
		&wf.CreatedBy,
		&wf.CreatedAt,
		&wf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if workflowType.Valid {
		wfType := models.WorkflowType(workflowType.String)
		wf.WorkflowType = &wfType
	}

	return &wf, nil
}

// GetByName retrieves a workflow template by name
func (r *WorkflowRepository) GetByName(name string) (*models.Workflow, error) {
	query := `
		SELECT id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows
		WHERE name = $1 AND is_active = true
	`

	var wf models.Workflow
	var workflowType sql.NullString
	err := r.db.QueryRow(query, name).Scan(
		&wf.ID,
		&wf.Name,
		&wf.Description,
		&workflowType,
		&wf.IsActive,
		&wf.CreatedBy,
		&wf.CreatedAt,
		&wf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if workflowType.Valid {
		wfType := models.WorkflowType(workflowType.String)
		wf.WorkflowType = &wfType
	}

	return &wf, nil
}

// GetByType retrieves a workflow template by its type
func (r *WorkflowRepository) GetByType(workflowType models.WorkflowType) (*models.Workflow, error) {
	query := `
		SELECT id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows
		WHERE workflow_type = $1 AND is_active = true
		LIMIT 1
	`

	var wf models.Workflow
	var wfTypeStr sql.NullString
	err := r.db.QueryRow(query, string(workflowType)).Scan(
		&wf.ID,
		&wf.Name,
		&wf.Description,
		&wfTypeStr,
		&wf.IsActive,
		&wf.CreatedBy,
		&wf.CreatedAt,
		&wf.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if wfTypeStr.Valid {
		wfType := models.WorkflowType(wfTypeStr.String)
		wf.WorkflowType = &wfType
	}

	return &wf, nil
}

// GetAllActive retrieves all active workflow templates
func (r *WorkflowRepository) GetAllActive() ([]models.Workflow, error) {
	query := `
		SELECT id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows
		WHERE is_active = true
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.Workflow
	for rows.Next() {
		var wf models.Workflow
		var workflowType sql.NullString
		err := rows.Scan(
			&wf.ID,
			&wf.Name,
			&wf.Description,
			&workflowType,
			&wf.IsActive,
			&wf.CreatedBy,
			&wf.CreatedAt,
			&wf.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if workflowType.Valid {
			wfType := models.WorkflowType(workflowType.String)
			wf.WorkflowType = &wfType
		}

		workflows = append(workflows, wf)
	}

	return workflows, nil
}

// GetAllActiveWithCounts retrieves all active workflow templates with step and transition counts
func (r *WorkflowRepository) GetAllActiveWithCounts() ([]models.WorkflowWithCounts, error) {
	query := `
		SELECT
			w.id,
			w.name,
			w.description,
			w.workflow_type,
			w.is_active,
			w.created_by,
			w.created_at,
			w.updated_at,
			COALESCE(COUNT(DISTINCT ws.id), 0) as step_count,
			COALESCE(COUNT(DISTINCT wt.id), 0) as transition_count
		FROM workflows w
		LEFT JOIN workflow_steps ws ON w.id = ws.workflow_id
		LEFT JOIN workflow_transitions wt ON w.id = wt.workflow_id
		WHERE w.is_active = true
		GROUP BY w.id, w.name, w.description, w.workflow_type, w.is_active, w.created_by, w.created_at, w.updated_at
		ORDER BY w.name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.WorkflowWithCounts
	for rows.Next() {
		var wf models.WorkflowWithCounts
		var workflowType sql.NullString
		err := rows.Scan(
			&wf.ID,
			&wf.Name,
			&wf.Description,
			&workflowType,
			&wf.IsActive,
			&wf.CreatedBy,
			&wf.CreatedAt,
			&wf.UpdatedAt,
			&wf.StepCount,
			&wf.TransitionCount,
		)
		if err != nil {
			return nil, err
		}

		if workflowType.Valid {
			wfType := models.WorkflowType(workflowType.String)
			wf.WorkflowType = &wfType
		}

		workflows = append(workflows, wf)
	}

	return workflows, nil
}

// GetStepsByWorkflowID retrieves all steps for a workflow
func (r *WorkflowRepository) GetStepsByWorkflowID(workflowID string) ([]models.WorkflowStep, error) {
	query := `
		SELECT id, workflow_id, step_name, step_order, initial, final,
		       allowed_roles, requires_all_approvers, min_approvals, created_at
		FROM workflow_steps
		WHERE workflow_id = $1
		ORDER BY step_order
	`

	rows, err := r.db.Query(query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []models.WorkflowStep
	for rows.Next() {
		var step models.WorkflowStep
		var allowedRolesJSON []byte

		err := rows.Scan(
			&step.ID,
			&step.WorkflowID,
			&step.StepName,
			&step.StepOrder,
			&step.Initial,
			&step.Final,
			&allowedRolesJSON,
			&step.RequiresAllApprovers,
			&step.MinApprovals,
			&step.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse allowed_roles JSON
		if err := json.Unmarshal(allowedRolesJSON, &step.AllowedRoles); err != nil {
			return nil, fmt.Errorf("failed to parse allowed_roles: %w", err)
		}

		steps = append(steps, step)
	}

	return steps, nil
}

// GetStepByID retrieves a specific workflow step
func (r *WorkflowRepository) GetStepByID(stepID string) (*models.WorkflowStep, error) {
	query := `
		SELECT id, workflow_id, step_name, step_order, initial, final,
		       allowed_roles, requires_all_approvers, min_approvals, created_at
		FROM workflow_steps
		WHERE id = $1
	`

	var step models.WorkflowStep
	var allowedRolesJSON []byte

	err := r.db.QueryRow(query, stepID).Scan(
		&step.ID,
		&step.WorkflowID,
		&step.StepName,
		&step.StepOrder,
		&step.Initial,
		&step.Final,
		&allowedRolesJSON,
		&step.RequiresAllApprovers,
		&step.MinApprovals,
		&step.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse allowed_roles JSON
	if err := json.Unmarshal(allowedRolesJSON, &step.AllowedRoles); err != nil {
		return nil, fmt.Errorf("failed to parse allowed_roles: %w", err)
	}

	return &step, nil
}

// GetInitialStep retrieves the initial step of a workflow
func (r *WorkflowRepository) GetInitialStep(workflowID string) (*models.WorkflowStep, error) {
	query := `
		SELECT id, workflow_id, step_name, step_order, initial, final,
		       allowed_roles, requires_all_approvers, min_approvals, created_at
		FROM workflow_steps
		WHERE workflow_id = $1 AND initial = true
		LIMIT 1
	`

	var step models.WorkflowStep
	var allowedRolesJSON []byte

	err := r.db.QueryRow(query, workflowID).Scan(
		&step.ID,
		&step.WorkflowID,
		&step.StepName,
		&step.StepOrder,
		&step.Initial,
		&step.Final,
		&allowedRolesJSON,
		&step.RequiresAllApprovers,
		&step.MinApprovals,
		&step.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Parse allowed_roles JSON
	if err := json.Unmarshal(allowedRolesJSON, &step.AllowedRoles); err != nil {
		return nil, fmt.Errorf("failed to parse allowed_roles: %w", err)
	}

	return &step, nil
}

// GetTransitionsByWorkflowID retrieves all transitions for a workflow
func (r *WorkflowRepository) GetTransitionsByWorkflowID(workflowID string) ([]models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       condition_type, condition_value, created_at
		FROM workflow_transitions
		WHERE workflow_id = $1
	`

	rows, err := r.db.Query(query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []models.WorkflowTransition
	for rows.Next() {
		var tr models.WorkflowTransition
		err := rows.Scan(
			&tr.ID,
			&tr.WorkflowID,
			&tr.FromStepID,
			&tr.ToStepID,
			&tr.ActionName,
			&tr.ConditionType,
			&tr.ConditionValue,
			&tr.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transitions = append(transitions, tr)
	}

	return transitions, nil
}

// GetValidTransitions retrieves valid transitions from a specific step
func (r *WorkflowRepository) GetValidTransitions(fromStepID string) ([]models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       condition_type, condition_value, created_at
		FROM workflow_transitions
		WHERE from_step_id = $1
	`

	rows, err := r.db.Query(query, fromStepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []models.WorkflowTransition
	for rows.Next() {
		var tr models.WorkflowTransition
		err := rows.Scan(
			&tr.ID,
			&tr.WorkflowID,
			&tr.FromStepID,
			&tr.ToStepID,
			&tr.ActionName,
			&tr.ConditionType,
			&tr.ConditionValue,
			&tr.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transitions = append(transitions, tr)
	}

	return transitions, nil
}

// GetTransitionByID retrieves a specific workflow transition by ID
func (r *WorkflowRepository) GetTransitionByID(id string) (*models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       condition_type, condition_value, created_at
		FROM workflow_transitions
		WHERE id = $1
	`

	var tr models.WorkflowTransition
	err := r.db.QueryRow(query, id).Scan(
		&tr.ID,
		&tr.WorkflowID,
		&tr.FromStepID,
		&tr.ToStepID,
		&tr.ActionName,
		&tr.ConditionType,
		&tr.ConditionValue,
		&tr.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &tr, nil
}

// GetTransitionByAction retrieves a specific transition based on action
func (r *WorkflowRepository) GetTransitionByAction(fromStepID, action string) (*models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       condition_type, condition_value, created_at
		FROM workflow_transitions
		WHERE from_step_id = $1 AND action_name = $2
		LIMIT 1
	`

	var tr models.WorkflowTransition
	err := r.db.QueryRow(query, fromStepID, action).Scan(
		&tr.ID,
		&tr.WorkflowID,
		&tr.FromStepID,
		&tr.ToStepID,
		&tr.ActionName,
		&tr.ConditionType,
		&tr.ConditionValue,
		&tr.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &tr, nil
}

// Create creates a new workflow template
func (r *WorkflowRepository) Create(workflow *models.Workflow) error {
	query := `
		INSERT INTO workflows (id, name, description, workflow_type, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	workflow.ID = uuid.New().String()

	var workflowType *string
	if workflow.WorkflowType != nil {
		typeStr := string(*workflow.WorkflowType)
		workflowType = &typeStr
	}

	return r.db.QueryRow(
		query,
		workflow.ID,
		workflow.Name,
		workflow.Description,
		workflowType,
		workflow.IsActive,
		workflow.CreatedBy,
	).Scan(&workflow.CreatedAt, &workflow.UpdatedAt)
}

// CreateStep creates a new workflow step
func (r *WorkflowRepository) CreateStep(step *models.WorkflowStep) error {
	query := `
		INSERT INTO workflow_steps (
			id, workflow_id, step_name, step_order, initial, final,
			allowed_roles, requires_all_approvers, min_approvals
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`

	step.ID = uuid.New().String()

	// Marshal allowed_roles to JSON
	allowedRolesJSON, err := json.Marshal(step.AllowedRoles)
	if err != nil {
		return fmt.Errorf("failed to marshal allowed_roles: %w", err)
	}

	return r.db.QueryRow(
		query,
		step.ID,
		step.WorkflowID,
		step.StepName,
		step.StepOrder,
		step.Initial,
		step.Final,
		allowedRolesJSON,
		step.RequiresAllApprovers,
		step.MinApprovals,
	).Scan(&step.CreatedAt)
}

// CreateTransition creates a new workflow transition
func (r *WorkflowRepository) CreateTransition(transition *models.WorkflowTransition) error {
	query := `
		INSERT INTO workflow_transitions (
			id, workflow_id, from_step_id, to_step_id, action_name,
			condition_type, condition_value
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	transition.ID = uuid.New().String()

	return r.db.QueryRow(
		query,
		transition.ID,
		transition.WorkflowID,
		transition.FromStepID,
		transition.ToStepID,
		transition.ActionName,
		transition.ConditionType,
		transition.ConditionValue,
	).Scan(&transition.CreatedAt)
}

// Update updates an existing workflow template
func (r *WorkflowRepository) Update(workflow *models.Workflow) error {
	query := `
		UPDATE workflows
		SET name = $1, description = $2, workflow_type = $3, is_active = $4
		WHERE id = $5
		RETURNING updated_at
	`

	var workflowType *string
	if workflow.WorkflowType != nil {
		typeStr := string(*workflow.WorkflowType)
		workflowType = &typeStr
	}

	return r.db.QueryRow(
		query,
		workflow.Name,
		workflow.Description,
		workflowType,
		workflow.IsActive,
		workflow.ID,
	).Scan(&workflow.UpdatedAt)
}

// UpdateStep updates an existing workflow step
func (r *WorkflowRepository) UpdateStep(step *models.WorkflowStep) error {
	query := `
		UPDATE workflow_steps
		SET step_name = $1, step_order = $2, initial = $3, final = $4,
		    allowed_roles = $5, requires_all_approvers = $6, min_approvals = $7
		WHERE id = $8
	`

	// Marshal allowed_roles to JSON
	allowedRolesJSON, err := json.Marshal(step.AllowedRoles)
	if err != nil {
		return fmt.Errorf("failed to marshal allowed_roles: %w", err)
	}

	result, err := r.db.Exec(
		query,
		step.StepName,
		step.StepOrder,
		step.Initial,
		step.Final,
		allowedRolesJSON,
		step.RequiresAllApprovers,
		step.MinApprovals,
		step.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow step not found")
	}

	return nil
}

// UpdateTransition updates an existing workflow transition
func (r *WorkflowRepository) UpdateTransition(transition *models.WorkflowTransition) error {
	query := `
		UPDATE workflow_transitions
		SET action_name = $1, condition_type = $2, condition_value = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(
		query,
		transition.ActionName,
		transition.ConditionType,
		transition.ConditionValue,
		transition.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow transition not found")
	}

	return nil
}

// Deactivate deactivates a workflow template by setting is_active to false
func (r *WorkflowRepository) Deactivate(id string) error {
	query := `
		UPDATE workflows
		SET is_active = false
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow not found")
	}

	return nil
}

// Delete permanently deletes a workflow template
func (r *WorkflowRepository) Delete(id string) error {
	query := `DELETE FROM workflows WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow not found")
	}

	return nil
}

// DeleteStep deletes a workflow step
func (r *WorkflowRepository) DeleteStep(id string) error {
	query := `DELETE FROM workflow_steps WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow step not found")
	}

	return nil
}

// DeleteTransition deletes a workflow transition
func (r *WorkflowRepository) DeleteTransition(id string) error {
	query := `DELETE FROM workflow_transitions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("workflow transition not found")
	}

	return nil
}

// GetFirstActionStep gets the first actionable step in a workflow (the step after submission)
// This is the step you transition to from the initial step via "submit" action
// For example: Initial Step (Draft) -> [submit] -> First Action Step (Pending Review)
// This combines the logic of: get initial step -> find submit transition -> get next step
// Much more efficient than doing 3 separate queries
func (r *WorkflowRepository) GetFirstActionStep(workflowID string) (*models.WorkflowStep, string, error) {
	query := `
		SELECT
			ws.id, ws.workflow_id, ws.step_name, ws.step_order, ws.initial, ws.final,
			ws.allowed_roles, ws.requires_all_approvers, ws.min_approvals, ws.created_at,
			initial.id as initial_step_id
		FROM workflow_steps ws
		INNER JOIN workflow_transitions wt ON ws.id = wt.to_step_id
		INNER JOIN workflow_steps initial ON wt.from_step_id = initial.id
		WHERE ws.workflow_id = $1
		  AND initial.initial = true
		  AND wt.action_name = 'submit'
		LIMIT 1
	`

	var step models.WorkflowStep
	var initialStepID string
	var allowedRolesJSON []byte

	err := r.db.QueryRow(query, workflowID).Scan(
		&step.ID,
		&step.WorkflowID,
		&step.StepName,
		&step.StepOrder,
		&step.Initial,
		&step.Final,
		&allowedRolesJSON,
		&step.RequiresAllApprovers,
		&step.MinApprovals,
		&step.CreatedAt,
		&initialStepID,
	)

	if err != nil {
		return nil, "", err
	}

	// Parse allowed_roles JSON
	if err := json.Unmarshal(allowedRolesJSON, &step.AllowedRoles); err != nil {
		return nil, "", fmt.Errorf("failed to parse allowed_roles: %w", err)
	}

	return &step, initialStepID, nil
}
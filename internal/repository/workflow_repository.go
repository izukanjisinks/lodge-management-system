package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

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

// GetByID retrieves a workflow template by ID, scoped to org.
func (r *WorkflowRepository) GetByID(id, orgID string) (*models.Workflow, error) {
	query := `
		SELECT id, org_id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows
		WHERE id = $1 AND org_id = $2
	`

	var wf models.Workflow
	var workflowType sql.NullString
	err := r.db.QueryRow(query, id, orgID).Scan(
		&wf.ID,
		&wf.OrgID,
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

// GetByName retrieves a workflow template by name, scoped to org.
func (r *WorkflowRepository) GetByName(name, orgID string) (*models.Workflow, error) {
	query := `
		SELECT id, org_id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows
		WHERE name = $1 AND is_active = true AND org_id = $2
	`

	var wf models.Workflow
	var workflowType sql.NullString
	err := r.db.QueryRow(query, name, orgID).Scan(
		&wf.ID,
		&wf.OrgID,
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

// GetByType retrieves a workflow template by its type, scoped to org.
func (r *WorkflowRepository) GetByType(workflowType models.WorkflowType, orgID string) (*models.Workflow, error) {
	var wf models.Workflow
	var wfTypeStr sql.NullString

	err := r.db.QueryRow(`
		SELECT id, org_id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows
		WHERE workflow_type = $1 AND is_active = true AND org_id = $2
		LIMIT 1`, string(workflowType), orgID).Scan(
		&wf.ID, &wf.OrgID, &wf.Name, &wf.Description, &wfTypeStr,
		&wf.IsActive, &wf.CreatedBy, &wf.CreatedAt, &wf.UpdatedAt,
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

// GetAllActive retrieves all active workflow templates, scoped to org.
func (r *WorkflowRepository) GetAllActive(orgID string) ([]models.Workflow, error) {
	rows, err := r.db.Query(`
		SELECT id, org_id, name, description, workflow_type, is_active, created_by, created_at, updated_at
		FROM workflows WHERE is_active = true AND org_id = $1 ORDER BY name`, orgID)
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
			&wf.OrgID,
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

// GetAllActiveWithCounts retrieves all active workflow templates with step and transition counts, scoped to org.
func (r *WorkflowRepository) GetAllActiveWithCounts(orgID string) ([]models.WorkflowWithCounts, error) {
	rows, err := r.db.Query(`
		SELECT
			w.id, w.name, w.description, w.workflow_type, w.is_active, w.created_by,
			w.created_at, w.updated_at,
			COALESCE(COUNT(DISTINCT ws.id), 0) as step_count,
			COALESCE(COUNT(DISTINCT wt.id), 0) as transition_count
		FROM workflows w
		LEFT JOIN workflow_steps ws ON w.id = ws.workflow_id
		LEFT JOIN workflow_transitions wt ON w.id = wt.workflow_id
		WHERE w.is_active = true AND w.org_id = $1
		GROUP BY w.id, w.name, w.description, w.workflow_type, w.is_active, w.created_by, w.created_at, w.updated_at
		ORDER BY w.name`, orgID)
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
		SELECT id, workflow_id, step_name, step_order, initial, final, created_at
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
		if err := rows.Scan(
			&step.ID,
			&step.WorkflowID,
			&step.StepName,
			&step.StepOrder,
			&step.Initial,
			&step.Final,
			&step.CreatedAt,
		); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	return steps, nil
}

// GetStepByID retrieves a specific workflow step
func (r *WorkflowRepository) GetStepByID(stepID string) (*models.WorkflowStep, error) {
	query := `
		SELECT id, workflow_id, step_name, step_order, initial, final, created_at
		FROM workflow_steps
		WHERE id = $1
	`

	var step models.WorkflowStep
	err := r.db.QueryRow(query, stepID).Scan(
		&step.ID,
		&step.WorkflowID,
		&step.StepName,
		&step.StepOrder,
		&step.Initial,
		&step.Final,
		&step.CreatedAt,
	)
	return &step, err
}

// GetInitialStep retrieves the initial step of a workflow
func (r *WorkflowRepository) GetInitialStep(workflowID string) (*models.WorkflowStep, error) {
	query := `
		SELECT id, workflow_id, step_name, step_order, initial, final, created_at
		FROM workflow_steps
		WHERE workflow_id = $1 AND initial = true
		LIMIT 1
	`

	var step models.WorkflowStep
	err := r.db.QueryRow(query, workflowID).Scan(
		&step.ID,
		&step.WorkflowID,
		&step.StepName,
		&step.StepOrder,
		&step.Initial,
		&step.Final,
		&step.CreatedAt,
	)
	return &step, err
}

type transitionScanner interface {
	Scan(dest ...interface{}) error
}

func scanTransition(row transitionScanner) (*models.WorkflowTransition, error) {
	var tr models.WorkflowTransition
	var allowedRolesJSON []byte
	err := row.Scan(
		&tr.ID, &tr.WorkflowID, &tr.FromStepID, &tr.ToStepID, &tr.ActionName,
		&allowedRolesJSON, &tr.ConditionType, &tr.ConditionValue, &tr.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(allowedRolesJSON, &tr.AllowedRoles); err != nil {
		return nil, fmt.Errorf("failed to parse allowed_roles: %w", err)
	}
	return &tr, nil
}

// GetTransitionsByWorkflowID retrieves all transitions for a workflow
func (r *WorkflowRepository) GetTransitionsByWorkflowID(workflowID string) ([]models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       allowed_roles, condition_type, condition_value, created_at
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
		tr, err := scanTransition(rows)
		if err != nil {
			return nil, err
		}
		transitions = append(transitions, *tr)
	}
	return transitions, nil
}

// GetValidTransitions retrieves valid transitions from a specific step
func (r *WorkflowRepository) GetValidTransitions(fromStepID string) ([]models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       allowed_roles, condition_type, condition_value, created_at
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
		tr, err := scanTransition(rows)
		if err != nil {
			return nil, err
		}
		transitions = append(transitions, *tr)
	}
	return transitions, nil
}

// GetTransitionByID retrieves a specific workflow transition by ID
func (r *WorkflowRepository) GetTransitionByID(id string) (*models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       allowed_roles, condition_type, condition_value, created_at
		FROM workflow_transitions
		WHERE id = $1
	`
	return scanTransition(r.db.QueryRow(query, id))
}

// GetTransitionByAction retrieves a specific transition based on action
func (r *WorkflowRepository) GetTransitionByAction(fromStepID, action string) (*models.WorkflowTransition, error) {
	query := `
		SELECT id, workflow_id, from_step_id, to_step_id, action_name,
		       allowed_roles, condition_type, condition_value, created_at
		FROM workflow_transitions
		WHERE from_step_id = $1 AND action_name = $2
		LIMIT 1
	`
	return scanTransition(r.db.QueryRow(query, fromStepID, action))
}

// Create creates a new workflow template.
func (r *WorkflowRepository) Create(workflow *models.Workflow) error {
	query := `
		INSERT INTO workflows (id, org_id, name, description, workflow_type, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
		workflow.OrgID,
		workflow.Name,
		workflow.Description,
		workflowType,
		workflow.IsActive,
		workflow.CreatedBy,
	).Scan(&workflow.CreatedAt, &workflow.UpdatedAt)
}

// CreateStep creates a new workflow step
func (r *WorkflowRepository) CreateStep(step *models.WorkflowStep) error {
	step.ID = uuid.New().String()
	return r.db.QueryRow(`
		INSERT INTO workflow_steps (id, workflow_id, step_name, step_order, initial, final)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`,
		step.ID, step.WorkflowID, step.StepName, step.StepOrder, step.Initial, step.Final,
	).Scan(&step.CreatedAt)
}

// CreateTransition creates a new workflow transition
func (r *WorkflowRepository) CreateTransition(transition *models.WorkflowTransition) error {
	transition.ID = uuid.New().String()

	allowedRolesJSON, err := json.Marshal(transition.AllowedRoles)
	if err != nil {
		return fmt.Errorf("failed to marshal allowed_roles: %w", err)
	}

	return r.db.QueryRow(`
		INSERT INTO workflow_transitions (
			id, workflow_id, from_step_id, to_step_id, action_name,
			allowed_roles, condition_type, condition_value
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at`,
		transition.ID, transition.WorkflowID, transition.FromStepID, transition.ToStepID,
		transition.ActionName, allowedRolesJSON, transition.ConditionType, transition.ConditionValue,
	).Scan(&transition.CreatedAt)
}

// Update updates an existing workflow template, scoped to org.
func (r *WorkflowRepository) Update(workflow *models.Workflow) error {
	query := `
		UPDATE workflows
		SET name = $1, description = $2, workflow_type = $3, is_active = $4
		WHERE id = $5 AND org_id = $6
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
		workflow.OrgID,
	).Scan(&workflow.UpdatedAt)
}

// UpdateStep updates an existing workflow step
func (r *WorkflowRepository) UpdateStep(step *models.WorkflowStep) error {
	result, err := r.db.Exec(`
		UPDATE workflow_steps
		SET step_name = $1, step_order = $2, initial = $3, final = $4
		WHERE id = $5`,
		step.StepName, step.StepOrder, step.Initial, step.Final, step.ID,
	)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("workflow step not found")
	}
	return nil
}

// UpdateTransition updates an existing workflow transition
func (r *WorkflowRepository) UpdateTransition(transition *models.WorkflowTransition) error {
	allowedRolesJSON, err := json.Marshal(transition.AllowedRoles)
	if err != nil {
		return fmt.Errorf("failed to marshal allowed_roles: %w", err)
	}

	result, err := r.db.Exec(`
		UPDATE workflow_transitions
		SET action_name = $1, allowed_roles = $2, condition_type = $3, condition_value = $4
		WHERE id = $5`,
		transition.ActionName, allowedRolesJSON, transition.ConditionType, transition.ConditionValue, transition.ID,
	)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("workflow transition not found")
	}
	return nil
}

// Deactivate deactivates a workflow template, scoped to org.
func (r *WorkflowRepository) Deactivate(id, orgID string) error {
	query := `
		UPDATE workflows
		SET is_active = false
		WHERE id = $1 AND org_id = $2
	`

	result, err := r.db.Exec(query, id, orgID)
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

// Delete permanently deletes a workflow template, scoped to org.
func (r *WorkflowRepository) Delete(id, orgID string) error {
	query := `DELETE FROM workflows WHERE id = $1 AND org_id = $2`

	result, err := r.db.Exec(query, id, orgID)
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

// GetFirstActionStep gets the first actionable step and the transition leading to it
// (from the initial step). Returns the step, the transition (with allowed_roles), and
// the initial step's ID.
func (r *WorkflowRepository) GetFirstActionStep(workflowID string) (*models.WorkflowStep, *models.WorkflowTransition, string, error) {
	var step models.WorkflowStep
	var transition models.WorkflowTransition
	var initialStepID string
	var allowedRolesJSON []byte

	err := r.db.QueryRow(`
		SELECT ws.id, ws.workflow_id, ws.step_name, ws.step_order, ws.initial, ws.final,
		       ws.created_at,
		       wt.id, wt.workflow_id, wt.from_step_id, wt.to_step_id,
		       wt.action_name, wt.allowed_roles, wt.condition_type, wt.condition_value, wt.created_at,
		       initial.id AS initial_step_id
		FROM workflow_steps ws
		INNER JOIN workflow_transitions wt ON ws.id = wt.to_step_id
		INNER JOIN workflow_steps initial ON wt.from_step_id = initial.id
		WHERE ws.workflow_id = $1
		  AND initial.initial = true
		ORDER BY ws.step_order ASC
		LIMIT 1`, workflowID,
	).Scan(
		&step.ID, &step.WorkflowID, &step.StepName, &step.StepOrder,
		&step.Initial, &step.Final, &step.CreatedAt,
		&transition.ID, &transition.WorkflowID, &transition.FromStepID, &transition.ToStepID,
		&transition.ActionName, &allowedRolesJSON, &transition.ConditionType, &transition.ConditionValue,
		&transition.CreatedAt, &initialStepID,
	)
	if err != nil {
		return nil, nil, "", err
	}

	if err := json.Unmarshal(allowedRolesJSON, &transition.AllowedRoles); err != nil {
		transition.AllowedRoles = []string{}
	}

	return &step, &transition, initialStepID, nil
}

// UpdateEntityStatus updates the status of the real entity that a workflow instance tracks.
// Called by ProcessAction when a workflow reaches a terminal state (approved or rejected).
func (r *WorkflowRepository) UpdateEntityStatus(entityType, entityID, status, orgID string) error {
	var query string

	switch entityType {
	case "booking":
		query = `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2 AND org_id = $3`
	default:
		return fmt.Errorf("unsupported entity type: %s", entityType)
	}

	res, err := r.db.Exec(query, status, entityID, orgID)
	if err != nil {
		return fmt.Errorf("failed to update %s status: %w", entityType, err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("%s not found: %s", entityType, entityID)
	}
	return nil
}
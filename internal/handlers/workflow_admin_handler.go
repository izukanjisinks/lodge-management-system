package handlers

import (
	"encoding/json"
	"net/http"

	"hr-system/internal/middleware"
	"hr-system/internal/models"
	"hr-system/internal/repository"
)

type WorkflowAdminHandler struct {
	workflowRepo *repository.WorkflowRepository
}

func NewWorkflowAdminHandler(workflowRepo *repository.WorkflowRepository) *WorkflowAdminHandler {
	return &WorkflowAdminHandler{
		workflowRepo: workflowRepo,
	}
}

// ========== Workflow Template Management ==========

// GetAllWorkflows retrieves all active workflow templates with counts
func (h *WorkflowAdminHandler) GetAllWorkflows(w http.ResponseWriter, r *http.Request) {
	workflows, err := h.workflowRepo.GetAllActiveWithCounts()
	if err != nil {
		http.Error(w, "Failed to retrieve workflows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"workflows": workflows,
		"count":     len(workflows),
	})
}

// GetWorkflowByID retrieves a specific workflow template
func (h *WorkflowAdminHandler) GetWorkflowByID(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("id")
	if workflowID == "" {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	workflow, err := h.workflowRepo.GetByID(workflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// CreateWorkflow creates a new workflow template
func (h *WorkflowAdminHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if workflow.Name == "" {
		http.Error(w, "Workflow name is required", http.StatusBadRequest)
		return
	}

	workflow.CreatedBy = userID.String()
	workflow.IsActive = true

	if err := h.workflowRepo.Create(&workflow); err != nil {
		http.Error(w, "Failed to create workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Workflow created successfully",
		"workflow": workflow,
	})
}

// UpdateWorkflowRequest represents the request body for updating a workflow
type UpdateWorkflowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// UpdateWorkflow updates an existing workflow template
func (h *WorkflowAdminHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("id")
	if workflowID == "" {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	// Get existing workflow
	workflow, err := h.workflowRepo.GetByID(workflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	var req UpdateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields if provided
	if req.Name != "" {
		workflow.Name = req.Name
	}
	if req.Description != "" {
		workflow.Description = req.Description
	}
	if req.IsActive != nil {
		workflow.IsActive = *req.IsActive
	}

	// Update workflow
	if err := h.workflowRepo.Update(workflow); err != nil {
		http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Workflow updated successfully",
		"workflow": workflow,
	})
}

// DeactivateWorkflow deactivates a workflow template (sets is_active to false)
func (h *WorkflowAdminHandler) DeactivateWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("id")
	if workflowID == "" {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	if err := h.workflowRepo.Deactivate(workflowID); err != nil {
		http.Error(w, "Failed to deactivate workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Workflow deactivated successfully",
	})
}

// DeleteWorkflow permanently deletes a workflow template
func (h *WorkflowAdminHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("id")
	if workflowID == "" {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	if err := h.workflowRepo.Delete(workflowID); err != nil {
		http.Error(w, "Failed to delete workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Workflow permanently deleted",
	})
}

// ========== Workflow Steps Management ==========

// GetWorkflowSteps retrieves all steps for a workflow
func (h *WorkflowAdminHandler) GetWorkflowSteps(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("id")
	if workflowID == "" {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	steps, err := h.workflowRepo.GetStepsByWorkflowID(workflowID)
	if err != nil {
		http.Error(w, "Failed to retrieve steps", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"steps": steps,
		"count": len(steps),
	})
}

// GetStepByID retrieves a specific workflow step
func (h *WorkflowAdminHandler) GetStepByID(w http.ResponseWriter, r *http.Request) {
	stepID := r.PathValue("step_id")
	if stepID == "" {
		http.Error(w, "Step ID is required", http.StatusBadRequest)
		return
	}

	step, err := h.workflowRepo.GetStepByID(stepID)
	if err != nil {
		http.Error(w, "Step not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(step)
}

// CreateWorkflowStepRequest represents the request body for creating a step
type CreateWorkflowStepRequest struct {
	WorkflowID           string   `json:"workflow_id"`
	StepName             string   `json:"step_name"`
	StepOrder            int      `json:"step_order"`
	Initial              bool     `json:"initial"`
	Final                bool     `json:"final"`
	AllowedRoles         []string `json:"allowed_roles"`
	RequiresAllApprovers bool     `json:"requires_all_approvers"`
	MinApprovals         int      `json:"min_approvals"`
}

// CreateWorkflowStep creates a new workflow step
func (h *WorkflowAdminHandler) CreateWorkflowStep(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkflowStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WorkflowID == "" || req.StepName == "" {
		http.Error(w, "Workflow ID and step name are required", http.StatusBadRequest)
		return
	}

	// Verify workflow exists
	_, err := h.workflowRepo.GetByID(req.WorkflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	// Create the step model
	step := &models.WorkflowStep{
		WorkflowID:           req.WorkflowID,
		StepName:             req.StepName,
		StepOrder:            req.StepOrder,
		Initial:              req.Initial,
		Final:                req.Final,
		AllowedRoles:         req.AllowedRoles,
		RequiresAllApprovers: req.RequiresAllApprovers,
		MinApprovals:         req.MinApprovals,
	}

	// Create the step
	if err := h.workflowRepo.CreateStep(step); err != nil {
		http.Error(w, "Failed to create workflow step", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Workflow step created successfully",
		"step":    step,
	})
}

// UpdateWorkflowStepRequest represents the request body for updating a step
type UpdateWorkflowStepRequest struct {
	StepName             *string   `json:"step_name"`
	StepOrder            *int      `json:"step_order"`
	Initial              *bool     `json:"initial"`
	Final                *bool     `json:"final"`
	AllowedRoles         *[]string `json:"allowed_roles"`
	RequiresAllApprovers *bool     `json:"requires_all_approvers"`
	MinApprovals         *int      `json:"min_approvals"`
}

// UpdateWorkflowStep updates an existing workflow step
func (h *WorkflowAdminHandler) UpdateWorkflowStep(w http.ResponseWriter, r *http.Request) {
	stepID := r.PathValue("step_id")
	if stepID == "" {
		http.Error(w, "Step ID is required", http.StatusBadRequest)
		return
	}

	// Get existing step
	step, err := h.workflowRepo.GetStepByID(stepID)
	if err != nil {
		http.Error(w, "Step not found", http.StatusNotFound)
		return
	}

	var req UpdateWorkflowStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields if provided
	if req.StepName != nil {
		step.StepName = *req.StepName
	}
	if req.StepOrder != nil {
		step.StepOrder = *req.StepOrder
	}
	if req.Initial != nil {
		step.Initial = *req.Initial
	}
	if req.Final != nil {
		step.Final = *req.Final
	}
	if req.AllowedRoles != nil {
		step.AllowedRoles = *req.AllowedRoles
	}
	if req.RequiresAllApprovers != nil {
		step.RequiresAllApprovers = *req.RequiresAllApprovers
	}
	if req.MinApprovals != nil {
		step.MinApprovals = *req.MinApprovals
	}

	// Update step
	if err := h.workflowRepo.UpdateStep(step); err != nil {
		http.Error(w, "Failed to update workflow step", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Workflow step updated successfully",
		"step":    step,
	})
}

// DeleteWorkflowStep deletes a workflow step
func (h *WorkflowAdminHandler) DeleteWorkflowStep(w http.ResponseWriter, r *http.Request) {
	stepID := r.PathValue("step_id")
	if stepID == "" {
		http.Error(w, "Step ID is required", http.StatusBadRequest)
		return
	}

	if err := h.workflowRepo.DeleteStep(stepID); err != nil {
		http.Error(w, "Failed to delete workflow step", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Workflow step deleted successfully",
	})
}

// ========== Workflow Transitions Management ==========

// GetWorkflowTransitions retrieves all transitions for a workflow
func (h *WorkflowAdminHandler) GetWorkflowTransitions(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("id")
	if workflowID == "" {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	transitions, err := h.workflowRepo.GetTransitionsByWorkflowID(workflowID)
	if err != nil {
		http.Error(w, "Failed to retrieve transitions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transitions": transitions,
		"count":       len(transitions),
	})
}

// GetValidTransitions retrieves valid transitions from a specific step
func (h *WorkflowAdminHandler) GetValidTransitions(w http.ResponseWriter, r *http.Request) {
	stepID := r.PathValue("step_id")
	if stepID == "" {
		http.Error(w, "Step ID is required", http.StatusBadRequest)
		return
	}

	transitions, err := h.workflowRepo.GetValidTransitions(stepID)
	if err != nil {
		http.Error(w, "Failed to retrieve transitions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transitions": transitions,
		"count":       len(transitions),
	})
}

// CreateWorkflowTransitionRequest represents the request body for creating a transition
type CreateWorkflowTransitionRequest struct {
	WorkflowID     string  `json:"workflow_id"`
	FromStepID     string  `json:"from_step_id"`
	ToStepID       string  `json:"to_step_id"`
	ActionName     string  `json:"action_name"`
	ConditionType  *string `json:"condition_type"`
	ConditionValue *string `json:"condition_value"`
}

// CreateWorkflowTransition creates a new workflow transition
func (h *WorkflowAdminHandler) CreateWorkflowTransition(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkflowTransitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WorkflowID == "" || req.FromStepID == "" || req.ToStepID == "" || req.ActionName == "" {
		http.Error(w, "Workflow ID, from step, to step, and action name are required", http.StatusBadRequest)
		return
	}

	// Verify workflow exists
	_, err := h.workflowRepo.GetByID(req.WorkflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	// Verify from step exists
	_, err = h.workflowRepo.GetStepByID(req.FromStepID)
	if err != nil {
		http.Error(w, "From step not found", http.StatusNotFound)
		return
	}

	// Verify to step exists
	_, err = h.workflowRepo.GetStepByID(req.ToStepID)
	if err != nil {
		http.Error(w, "To step not found", http.StatusNotFound)
		return
	}

	// Create the transition model
	transition := &models.WorkflowTransition{
		WorkflowID:     req.WorkflowID,
		FromStepID:     req.FromStepID,
		ToStepID:       req.ToStepID,
		ActionName:     req.ActionName,
		ConditionType:  req.ConditionType,
		ConditionValue: req.ConditionValue,
	}

	// Create the transition
	if err := h.workflowRepo.CreateTransition(transition); err != nil {
		http.Error(w, "Failed to create workflow transition", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Workflow transition created successfully",
		"transition": transition,
	})
}

// UpdateWorkflowTransitionRequest represents the request body for updating a transition
type UpdateWorkflowTransitionRequest struct {
	ActionName     *string `json:"action_name"`
	ConditionType  *string `json:"condition_type"`
	ConditionValue *string `json:"condition_value"`
}

// UpdateWorkflowTransition updates an existing workflow transition
func (h *WorkflowAdminHandler) UpdateWorkflowTransition(w http.ResponseWriter, r *http.Request) {
	transitionID := r.PathValue("transition_id")
	if transitionID == "" {
		http.Error(w, "Transition ID is required", http.StatusBadRequest)
		return
	}

	// Get existing transition
	transition, err := h.workflowRepo.GetTransitionByID(transitionID)
	if err != nil {
		http.Error(w, "Transition not found", http.StatusNotFound)
		return
	}

	var req UpdateWorkflowTransitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update fields if provided
	if req.ActionName != nil {
		transition.ActionName = *req.ActionName
	}
	if req.ConditionType != nil {
		transition.ConditionType = req.ConditionType
	}
	if req.ConditionValue != nil {
		transition.ConditionValue = req.ConditionValue
	}

	// Update transition
	if err := h.workflowRepo.UpdateTransition(transition); err != nil {
		http.Error(w, "Failed to update workflow transition", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Workflow transition updated successfully",
		"transition": transition,
	})
}

// DeleteWorkflowTransition deletes a workflow transition
func (h *WorkflowAdminHandler) DeleteWorkflowTransition(w http.ResponseWriter, r *http.Request) {
	transitionID := r.PathValue("transition_id")
	if transitionID == "" {
		http.Error(w, "Transition ID is required", http.StatusBadRequest)
		return
	}

	if err := h.workflowRepo.DeleteTransition(transitionID); err != nil {
		http.Error(w, "Failed to delete workflow transition", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Workflow transition deleted successfully",
	})
}

// ========== Workflow Structure Overview ==========

// GetWorkflowStructure retrieves complete workflow structure (steps + transitions)
func (h *WorkflowAdminHandler) GetWorkflowStructure(w http.ResponseWriter, r *http.Request) {
	workflowID := r.PathValue("id")
	if workflowID == "" {
		http.Error(w, "Workflow ID is required", http.StatusBadRequest)
		return
	}

	// Get workflow
	workflow, err := h.workflowRepo.GetByID(workflowID)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	// Get steps
	steps, err := h.workflowRepo.GetStepsByWorkflowID(workflowID)
	if err != nil {
		http.Error(w, "Failed to retrieve steps", http.StatusInternalServerError)
		return
	}

	// Get transitions
	transitions, err := h.workflowRepo.GetTransitionsByWorkflowID(workflowID)
	if err != nil {
		http.Error(w, "Failed to retrieve transitions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"workflow":    workflow,
		"steps":       steps,
		"transitions": transitions,
	})
}

// GetAvailableWorkflowTypes returns all available workflow type constants
// This endpoint helps frontend/API consumers know what workflow types can be created
func (h *WorkflowAdminHandler) GetAvailableWorkflowTypes(w http.ResponseWriter, r *http.Request) {
	workflowTypes := models.GetWorkflowTypeInfo()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"workflow_types": workflowTypes,
		"count":          len(workflowTypes),
	})
}
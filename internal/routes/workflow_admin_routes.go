package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

// RegisterWorkflowAdminRoutes registers all workflow administration routes
// These routes are for managing workflow templates, steps, and transitions
// Access should be restricted to super_admin and hr_manager roles
func RegisterWorkflowAdminRoutes(h *handlers.WorkflowAdminHandler) {
	// ========== Workflow Type Information ==========
	// GET available workflow types (constants)
	http.HandleFunc("GET /api/v1/admin/workflow-types", withAuthAndRole(h.GetAvailableWorkflowTypes, models.RoleAdmin))

	// ========== Workflow Template Management ==========
	// POST create workflow
	http.HandleFunc("POST /api/v1/admin/workflows", withAuthAndRole(h.CreateWorkflow, models.RoleAdmin))

	// GET all workflows
	http.HandleFunc("GET /api/v1/admin/workflows", withAuthAndRole(h.GetAllWorkflows, models.RoleAdmin))

	// ========== Workflow Steps Management ==========
	// POST create step
	http.HandleFunc("POST /api/v1/admin/workflow-steps", withAuthAndRole(h.CreateWorkflowStep, models.RoleAdmin))

	// PUT update step
	http.HandleFunc("PUT /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.UpdateWorkflowStep, models.RoleAdmin))

	// DELETE step
	http.HandleFunc("DELETE /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.DeleteWorkflowStep, models.RoleAdmin))

	// GET specific step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.GetStepByID, models.RoleAdmin))

	// GET valid transitions from a step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}/transitions", withAuthAndRole(h.GetValidTransitions, models.RoleAdmin))

	// ========== Workflow Transitions Management ==========
	// POST create transition
	http.HandleFunc("POST /api/v1/admin/workflow-transitions", withAuthAndRole(h.CreateWorkflowTransition, models.RoleAdmin))

	// PUT update transition
	http.HandleFunc("PUT /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.UpdateWorkflowTransition, models.RoleAdmin))

	// DELETE transition
	http.HandleFunc("DELETE /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.DeleteWorkflowTransition, models.RoleAdmin))

	// ========== Workflow-specific routes (must come after step/transition routes) ==========
	// GET complete workflow structure (steps + transitions)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/structure", withAuthAndRole(h.GetWorkflowStructure, models.RoleAdmin))

	// GET all steps for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/steps", withAuthAndRole(h.GetWorkflowSteps, models.RoleAdmin))

	// GET all transitions for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/transitions", withAuthAndRole(h.GetWorkflowTransitions, models.RoleAdmin))

	// PUT update workflow
	http.HandleFunc("PUT /api/v1/admin/workflows/{id}", withAuthAndRole(h.UpdateWorkflow, models.RoleAdmin))

	// DELETE deactivate workflow (soft delete)
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}/deactivate", withAuthAndRole(h.DeactivateWorkflow, models.RoleAdmin))

	// DELETE workflow permanently
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}", withAuthAndRole(h.DeleteWorkflow, models.RoleAdmin))

	// GET specific workflow (must be last to avoid conflicts)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}", withAuthAndRole(h.GetWorkflowByID, models.RoleAdmin))
}

package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
)

// RegisterWorkflowAdminRoutes registers all workflow administration routes
// These routes are for managing workflow templates, steps, and transitions
// Access should be restricted to super_admin and hr_manager roles
func RegisterWorkflowAdminRoutes(h *handlers.WorkflowAdminHandler) {
	// ========== Workflow Type Information ==========
	// GET available workflow types (constants)
	http.HandleFunc("GET /api/v1/admin/workflow-types", withAuthAndRole(h.GetAvailableWorkflowTypes, "admin"))

	// ========== Workflow Template Management ==========
	// POST create workflow
	http.HandleFunc("POST /api/v1/admin/workflows", withAuthAndRole(h.CreateWorkflow, "admin"))

	// GET all workflows
	http.HandleFunc("GET /api/v1/admin/workflows", withAuthAndRole(h.GetAllWorkflows, "admin"))

	// ========== Workflow Steps Management ==========
	// POST create step
	http.HandleFunc("POST /api/v1/admin/workflow-steps", withAuthAndRole(h.CreateWorkflowStep, "admin"))

	// PUT update step
	http.HandleFunc("PUT /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.UpdateWorkflowStep, "admin"))

	// DELETE step
	http.HandleFunc("DELETE /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.DeleteWorkflowStep, "admin"))

	// GET specific step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.GetStepByID, "admin"))

	// GET valid transitions from a step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}/transitions", withAuthAndRole(h.GetValidTransitions, "admin"))

	// ========== Workflow Transitions Management ==========
	// POST create transition
	http.HandleFunc("POST /api/v1/admin/workflow-transitions", withAuthAndRole(h.CreateWorkflowTransition, "admin"))

	// PUT update transition
	http.HandleFunc("PUT /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.UpdateWorkflowTransition, "admin"))

	// DELETE transition
	http.HandleFunc("DELETE /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.DeleteWorkflowTransition, "admin"))

	// ========== Workflow-specific routes (must come after step/transition routes) ==========
	// GET complete workflow structure (steps + transitions)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/structure", withAuthAndRole(h.GetWorkflowStructure, "admin"))

	// GET all steps for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/steps", withAuthAndRole(h.GetWorkflowSteps, "admin"))

	// GET all transitions for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/transitions", withAuthAndRole(h.GetWorkflowTransitions, "admin"))

	// PUT update workflow
	http.HandleFunc("PUT /api/v1/admin/workflows/{id}", withAuthAndRole(h.UpdateWorkflow, "admin"))

	// DELETE deactivate workflow (soft delete)
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}/deactivate", withAuthAndRole(h.DeactivateWorkflow, "admin"))

	// DELETE workflow permanently
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}", withAuthAndRole(h.DeleteWorkflow, "admin"))

	// GET specific workflow (must be last to avoid conflicts)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}", withAuthAndRole(h.GetWorkflowByID, "admin"))
}

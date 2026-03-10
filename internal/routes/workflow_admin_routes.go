package routes

import (
	"net/http"

	"hr-system/internal/handlers"
)

// RegisterWorkflowAdminRoutes registers all workflow administration routes
// These routes are for managing workflow templates, steps, and transitions
// Access should be restricted to super_admin and hr_manager roles
func RegisterWorkflowAdminRoutes(h *handlers.WorkflowAdminHandler) {
	// ========== Workflow Type Information ==========
	// GET available workflow types (constants)
	http.HandleFunc("GET /api/v1/admin/workflow-types", withAuthAndRole(h.GetAvailableWorkflowTypes, "super_admin", "hr_manager"))

	// ========== Workflow Template Management ==========
	// POST create workflow
	http.HandleFunc("POST /api/v1/admin/workflows", withAuthAndRole(h.CreateWorkflow, "super_admin", "hr_manager"))

	// GET all workflows
	http.HandleFunc("GET /api/v1/admin/workflows", withAuthAndRole(h.GetAllWorkflows, "super_admin", "hr_manager"))

	// ========== Workflow Steps Management ==========
	// POST create step
	http.HandleFunc("POST /api/v1/admin/workflow-steps", withAuthAndRole(h.CreateWorkflowStep, "super_admin", "hr_manager"))

	// PUT update step
	http.HandleFunc("PUT /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.UpdateWorkflowStep, "super_admin", "hr_manager"))

	// DELETE step
	http.HandleFunc("DELETE /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.DeleteWorkflowStep, "super_admin", "hr_manager"))

	// GET specific step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.GetStepByID, "super_admin", "hr_manager"))

	// GET valid transitions from a step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}/transitions", withAuthAndRole(h.GetValidTransitions, "super_admin", "hr_manager"))

	// ========== Workflow Transitions Management ==========
	// POST create transition
	http.HandleFunc("POST /api/v1/admin/workflow-transitions", withAuthAndRole(h.CreateWorkflowTransition, "super_admin", "hr_manager"))

	// PUT update transition
	http.HandleFunc("PUT /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.UpdateWorkflowTransition, "super_admin", "hr_manager"))

	// DELETE transition
	http.HandleFunc("DELETE /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.DeleteWorkflowTransition, "super_admin", "hr_manager"))

	// ========== Workflow-specific routes (must come after step/transition routes) ==========
	// GET complete workflow structure (steps + transitions)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/structure", withAuthAndRole(h.GetWorkflowStructure, "super_admin", "hr_manager"))

	// GET all steps for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/steps", withAuthAndRole(h.GetWorkflowSteps, "super_admin", "hr_manager"))

	// GET all transitions for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/transitions", withAuthAndRole(h.GetWorkflowTransitions, "super_admin", "hr_manager"))

	// PUT update workflow
	http.HandleFunc("PUT /api/v1/admin/workflows/{id}", withAuthAndRole(h.UpdateWorkflow, "super_admin", "hr_manager"))

	// DELETE deactivate workflow (soft delete)
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}/deactivate", withAuthAndRole(h.DeactivateWorkflow, "super_admin", "hr_manager"))

	// DELETE workflow permanently
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}", withAuthAndRole(h.DeleteWorkflow, "super_admin", "hr_manager"))

	// GET specific workflow (must be last to avoid conflicts)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}", withAuthAndRole(h.GetWorkflowByID, "super_admin", "hr_manager"))
}

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
	http.HandleFunc("GET /api/v1/admin/workflow-types", withAuthAndRole(h.GetAvailableWorkflowTypes, models.RoleAdmin, models.RoleBranchAdmin))

	// ========== Workflow Template Management ==========
	// POST create workflow
	http.HandleFunc("POST /api/v1/admin/workflows", withAuthAndRole(h.CreateWorkflow, models.RoleAdmin, models.RoleBranchAdmin))

	// GET all workflows
	http.HandleFunc("GET /api/v1/admin/workflows", withAuthAndRole(h.GetAllWorkflows, models.RoleAdmin, models.RoleBranchAdmin))

	// ========== Workflow Steps Management ==========
	// POST create step
	http.HandleFunc("POST /api/v1/admin/workflow-steps", withAuthAndRole(h.CreateWorkflowStep, models.RoleAdmin, models.RoleBranchAdmin))

	// PUT update step
	http.HandleFunc("PUT /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.UpdateWorkflowStep, models.RoleAdmin, models.RoleBranchAdmin))

	// DELETE step
	http.HandleFunc("DELETE /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.DeleteWorkflowStep, models.RoleAdmin, models.RoleBranchAdmin))

	// GET specific step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}", withAuthAndRole(h.GetStepByID, models.RoleAdmin, models.RoleBranchAdmin))

	// GET valid transitions from a step
	http.HandleFunc("GET /api/v1/admin/workflow-steps/{step_id}/transitions", withAuthAndRole(h.GetValidTransitions, models.RoleAdmin, models.RoleBranchAdmin))

	// ========== Workflow Transitions Management ==========
	// POST create transition
	http.HandleFunc("POST /api/v1/admin/workflow-transitions", withAuthAndRole(h.CreateWorkflowTransition, models.RoleAdmin, models.RoleBranchAdmin))

	// PUT update transition
	http.HandleFunc("PUT /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.UpdateWorkflowTransition, models.RoleAdmin, models.RoleBranchAdmin))

	// DELETE transition
	http.HandleFunc("DELETE /api/v1/admin/workflow-transitions/{transition_id}", withAuthAndRole(h.DeleteWorkflowTransition, models.RoleAdmin, models.RoleBranchAdmin))

	// ========== Workflow-specific routes (must come after step/transition routes) ==========
	// GET complete workflow structure (steps + transitions)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/structure", withAuthAndRole(h.GetWorkflowStructure, models.RoleAdmin, models.RoleBranchAdmin))

	// GET all steps for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/steps", withAuthAndRole(h.GetWorkflowSteps, models.RoleAdmin, models.RoleBranchAdmin))

	// GET all transitions for a workflow
	http.HandleFunc("GET /api/v1/admin/workflows/{id}/transitions", withAuthAndRole(h.GetWorkflowTransitions, models.RoleAdmin, models.RoleBranchAdmin))

	// PUT update workflow
	http.HandleFunc("PUT /api/v1/admin/workflows/{id}", withAuthAndRole(h.UpdateWorkflow, models.RoleAdmin, models.RoleBranchAdmin))

	// DELETE deactivate workflow (soft delete)
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}/deactivate", withAuthAndRole(h.DeactivateWorkflow, models.RoleAdmin, models.RoleBranchAdmin))

	// DELETE workflow permanently
	http.HandleFunc("DELETE /api/v1/admin/workflows/{id}", withAuthAndRole(h.DeleteWorkflow, models.RoleAdmin, models.RoleBranchAdmin))

	// GET specific workflow (must be last to avoid conflicts)
	http.HandleFunc("GET /api/v1/admin/workflows/{id}", withAuthAndRole(h.GetWorkflowByID, models.RoleAdmin, models.RoleBranchAdmin))
}

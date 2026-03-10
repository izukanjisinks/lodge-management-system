package routes

import (
	"hr-system/internal/handlers"
	"hr-system/internal/models"
	"net/http"
)

func RegisterWorkflowRoutes(h *handlers.WorkflowHandler) {
	// Get my tasks (all or filtered by status)
	// Query param: ?status=pending
	http.HandleFunc("GET /api/v1/workflow/my-tasks",
		withAuthAndRole(h.GetMyTasks, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	// Get only pending tasks
	http.HandleFunc("GET /api/v1/workflow/my-tasks/pending",
		withAuthAndRole(h.GetMyPendingTasks, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	// Get specific task details with full context
	http.HandleFunc("GET /api/v1/workflow/tasks/{id}",
		withAuthAndRole(h.GetTaskDetails, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	// Initiate a new workflow (typically called by other services, not directly by users)
	http.HandleFunc("POST /api/v1/workflow/instances",
		withAuthAndRole(h.InitiateWorkflow, models.RoleSuperAdmin, models.RoleHRManager))

	// Get workflow instance by task ID (e.g., leave request ID)
	// This needs to come before the {id} routes to avoid conflicts
	http.HandleFunc("GET /api/v1/workflow/task/{task_id}/instance",
		withAuthAndRole(h.GetInstanceByTaskID, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	// Process an action on a workflow instance
	// Body: {"action": "approve", "comments": "Looks good"}
	http.HandleFunc("POST /api/v1/workflow/instances/{id}/action",
		withAuthAndRole(h.ProcessAction, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))

	// Get workflow instance history
	http.HandleFunc("GET /api/v1/workflow/instances/{id}/history", withAuthAndRole(h.GetInstanceHistory, models.RoleSuperAdmin, models.RoleHRManager, models.RoleManager))
}

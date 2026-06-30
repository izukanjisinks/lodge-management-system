package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterWorkflowRoutes(h *handlers.WorkflowHandler) {
	http.HandleFunc("GET /api/v1/workflow/all-tasks",
		withAuthAndRole(h.GetAllOrgTasks, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("GET /api/v1/workflow/my-tasks",
		withAuthAndRole(h.GetMyTasks, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("GET /api/v1/workflow/my-tasks/pending",
		withAuthAndRole(h.GetMyPendingTasks, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("GET /api/v1/workflow/tasks/{id}",
		withAuthAndRole(h.GetTaskDetails, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("POST /api/v1/workflow/instances",
		withAuthAndRole(h.InitiateWorkflow, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("GET /api/v1/workflow/task/{task_id}/instance",
		withAuthAndRole(h.GetInstanceByTaskID, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("POST /api/v1/workflow/instances/{id}/action",
		withAuthAndRole(h.ProcessAction, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("GET /api/v1/workflow/instances/{id}/history",
		withAuthAndRole(h.GetInstanceHistory, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
}

package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterDashboardRoutes(h *handlers.DashboardHandler) {
	// Staff dashboard — admin, manager, receptionist
	http.HandleFunc("GET /api/v1/dashboard/stats",
		withAuthAndRole(h.StaffStats, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
}

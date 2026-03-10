package routes

import (
	"hr-system/internal/handlers"
	"hr-system/internal/models"
	"net/http"
)

func RegisterDashboardRoutes(h *handlers.DashboardHandler) {
	http.HandleFunc("GET /api/v1/hr/dashboard/me", withAuth(h.GetMyDashboard))
	http.HandleFunc("GET /api/v1/hr/dashboard/admin", withAuthAndRole(h.GetAdminDashboard, models.RoleSuperAdmin, models.RoleHRManager))
}

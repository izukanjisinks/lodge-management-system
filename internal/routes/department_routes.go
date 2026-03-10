package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterDepartmentRoutes(h *handlers.DepartmentHandler) {
	http.HandleFunc("GET /api/v1/hr/departments/tree",
		withAuth(h.GetTree))

	http.HandleFunc("GET /api/v1/hr/departments",
		withAuth(h.List))

	http.HandleFunc("POST /api/v1/hr/departments",
		withAuthAndRole(h.Create, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("GET /api/v1/hr/departments/{id}",
		withAuth(h.GetByID))

	http.HandleFunc("PUT /api/v1/hr/departments/{id}",
		withAuthAndRole(h.Update, models.RoleSuperAdmin, models.RoleHRManager))

	http.HandleFunc("DELETE /api/v1/hr/departments/{id}",
		withAuthAndRole(h.Delete, models.RoleSuperAdmin))
}
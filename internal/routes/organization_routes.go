package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterOrganizationRoutes(h *handlers.OrganizationHandler) {
	http.HandleFunc("GET /api/v1/organization",
		withAuthAndRole(h.Get, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("GET /api/v1/organization/{id}",
		withAuthAndRole(h.GetByID, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("PUT /api/v1/organization",
		withAuthAndRole(h.Update, models.RoleAdmin))
}

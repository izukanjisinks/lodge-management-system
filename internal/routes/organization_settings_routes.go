package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterOrganizationSettingsRoutes(h *handlers.OrganizationSettingsHandler) {
	http.HandleFunc("GET /api/v1/settings",
		withAuthAndRole(h.Get, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/settings",
		withAuthAndRole(h.Upsert, models.RoleAdmin))
}

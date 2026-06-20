package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterVenueRoutes(h *handlers.VenueHandler) {
	// Read — staff only (org-scoped via JWT)
	http.HandleFunc("GET /api/v1/venues",
		withAuth(h.List))

	http.HandleFunc("GET /api/v1/venues/{id}",
		withAuth(h.GetByID))

	// Write — admin and manager only
	http.HandleFunc("POST /api/v1/venues",
		withAuthAndRole(h.Create, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/venues/{id}",
		withAuthAndRole(h.Update, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/venues/{id}/images",
		withAuthAndRole(h.UpdateImages, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("PATCH /api/v1/venues/{id}/availability",
		withAuthAndRole(h.SetAvailability, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("DELETE /api/v1/venues/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))
}

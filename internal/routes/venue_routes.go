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

	// Write — branch admin and manager only (admin has no branch to scope creation to)
	http.HandleFunc("POST /api/v1/venues",
		withAuthAndRole(h.Create, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/venues/{id}",
		withAuthAndRole(h.Update, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/venues/{id}/images",
		withAuthAndRole(h.UpdateImages, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("PATCH /api/v1/venues/{id}/availability",
		withAuthAndRole(h.SetAvailability, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("DELETE /api/v1/venues/{id}",
		withAuthAndRole(h.Delete, models.RoleBranchAdmin, models.RoleManager))
}

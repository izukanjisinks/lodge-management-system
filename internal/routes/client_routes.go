package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterClientRoutes(h *handlers.ClientHandler) {
	// Individual clients — read: all staff; write: admin, manager, receptionist
	http.HandleFunc("GET /api/v1/clients/individual/lookup",
		withAuth(h.LookupIndividualByIDNumber))
	http.HandleFunc("GET /api/v1/clients/individual",
		withAuth(h.ListIndividual))

	http.HandleFunc("GET /api/v1/clients/individual/{id}",
		withAuth(h.GetIndividualByID))

	http.HandleFunc("POST /api/v1/clients/individual",
		withAuthAndRole(h.CreateIndividual, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("PUT /api/v1/clients/individual/{id}",
		withAuthAndRole(h.UpdateIndividual, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("DELETE /api/v1/clients/individual/{id}",
		withAuthAndRole(h.DeleteIndividual, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))

	// Corporate clients — same role split
	http.HandleFunc("GET /api/v1/clients/corporate/search",
		withAuth(h.SearchCorporate))
	http.HandleFunc("GET /api/v1/clients/corporate",
		withAuth(h.ListCorporate))

	http.HandleFunc("GET /api/v1/clients/corporate/{id}",
		withAuth(h.GetCorporateByID))

	http.HandleFunc("GET /api/v1/clients/corporate/{id}/bookings",
		withAuth(h.GetCorporateWithBookings))

	http.HandleFunc("PUT /api/v1/clients/corporate/{id}/documents",
		withAuthAndRole(h.UpsertDocuments, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("POST /api/v1/clients/corporate",
		withAuthAndRole(h.CreateCorporate, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("PUT /api/v1/clients/corporate/{id}",
		withAuthAndRole(h.UpdateCorporate, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("DELETE /api/v1/clients/corporate/{id}",
		withAuthAndRole(h.DeleteCorporate, models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))
}

package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterBookingRoutes(h *handlers.BookingHandler) {
	// Read — all authenticated staff
	http.HandleFunc("GET /api/v1/bookings",
		withAuth(h.List))

	http.HandleFunc("GET /api/v1/bookings/{id}",
		withAuth(h.GetByID))

	// Create/update — admin, manager, receptionist
	http.HandleFunc("POST /api/v1/bookings",
		withAuthAndRole(h.Create, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("PUT /api/v1/bookings/{id}",
		withAuthAndRole(h.Update, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	// Status transitions — admin and manager only
	http.HandleFunc("PATCH /api/v1/bookings/{id}/status",
		withAuthAndRole(h.UpdateStatus, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	// Clear overstayed flag — admin and manager only (set only by the nightly job)
	http.HandleFunc("PATCH /api/v1/bookings/{id}/clear-overstayed",
		withAuthAndRole(h.ClearOverstayed, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	// Delete — admin and manager only
	http.HandleFunc("DELETE /api/v1/bookings/{id}",
		withAuthAndRole(h.Delete, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
}

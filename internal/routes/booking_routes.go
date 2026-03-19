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
		withAuthAndRole(h.Create, models.RoleAdmin, models.RoleManager, models.RoleReceptionist))

	http.HandleFunc("PUT /api/v1/bookings/{id}",
		withAuthAndRole(h.Update, models.RoleAdmin, models.RoleManager, models.RoleReceptionist))

	// Status transitions — admin and manager only
	http.HandleFunc("PATCH /api/v1/bookings/{id}/status",
		withAuthAndRole(h.UpdateStatus, models.RoleAdmin, models.RoleManager))

	// Delete — admin and manager only
	http.HandleFunc("DELETE /api/v1/bookings/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin, models.RoleManager))
}

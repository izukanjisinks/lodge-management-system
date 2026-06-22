package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterIndividualBookingRequestRoutes(h *handlers.IndividualBookingRequestHandler) {
	// Web user — submit accommodation booking request (unified envelope from frontend).
	// URL path uses "guest" to match the frontend, but auth validates against web_users.
	http.HandleFunc("POST /api/v1/guest/bookings/accommodation", withWebUserAuth(h.SubmitAccommodation))
	// Web user — submit standalone event booking request (Flow B)
	http.HandleFunc("POST /api/v1/guest/bookings/event", withWebUserAuth(h.SubmitEvent))

	// Web user — submit and manage their own requests
	http.HandleFunc("POST /api/v1/web/bookings", withWebUserAuth(h.Submit))
	http.HandleFunc("GET /api/v1/web/bookings", withWebUserAuth(h.ListForWebUser))
	http.HandleFunc("GET /api/v1/web/bookings/{id}", withWebUserAuth(h.GetForWebUser))
	http.HandleFunc("PATCH /api/v1/web/bookings/{id}/cancel", withWebUserAuth(h.CancelForWebUser))

	// Backoffice — review and action requests
	http.HandleFunc("GET /api/v1/booking-requests/individual",
		withAuth(h.List))
	http.HandleFunc("GET /api/v1/booking-requests/individual/{id}",
		withAuth(h.GetByID))
	http.HandleFunc("PUT /api/v1/booking-requests/individual/{id}/approve",
		withAuthAndRole(h.Approve,
			models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/booking-requests/individual/{id}/reject",
		withAuthAndRole(h.Reject,
			models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/booking-requests/individual/{id}/cancel",
		withAuthAndRole(h.Cancel,
			models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
}

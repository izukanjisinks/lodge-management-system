package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterBookingRoutes(h *handlers.BookingHandler) {
	// ─── Web user confirmed bookings ───────────────────────────────────────────
	http.HandleFunc("GET /api/v1/web/my-bookings",
		withWebUserAuth(h.ListForWebUser))

	// ─── Bookings ──────────────────────────────────────────────────────────────
	http.HandleFunc("GET /api/v1/bookings",
		withAuth(h.List))
	http.HandleFunc("GET /api/v1/bookings/{id}",
		withAuth(h.GetByID))

	http.HandleFunc("POST /api/v1/bookings/individual",
		withAuthAndRole(h.CreateIndividual, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("POST /api/v1/booking-requests/{request_id}/materialise",
		withAuthAndRole(h.CreateFromRequest, models.RoleBranchAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/bookings/{id}/status",
		withAuthAndRole(h.UpdateStatus, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("PUT /api/v1/bookings/{id}/checkin",
		withAuthAndRole(h.CheckIn, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("PUT /api/v1/bookings/{id}/checkout",
		withAuthAndRole(h.CheckOut, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("DELETE /api/v1/bookings/{id}",
		withAuthAndRole(h.Cancel, models.RoleBranchAdmin, models.RoleManager))

	// ─── Room assignments ──────────────────────────────────────────────────────
	http.HandleFunc("GET /api/v1/bookings/{id}/assignments",
		withAuth(h.ListAssignments))
	http.HandleFunc("POST /api/v1/bookings/{id}/assignments",
		withAuthAndRole(h.AssignRoom, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("PUT /api/v1/bookings/{id}/assignments/{assign_id}",
		withAuthAndRole(h.UpdateAssignment, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("DELETE /api/v1/bookings/{id}/assignments/{assign_id}",
		withAuthAndRole(h.RemoveAssignment, models.RoleBranchAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/bookings/{id}/assignments/{assign_id}/checkin",
		withAuthAndRole(h.CheckInAssignment, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("PUT /api/v1/bookings/{id}/assignments/{assign_id}/checkout",
		withAuthAndRole(h.CheckOutAssignment, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))

	// ─── Attendees ─────────────────────────────────────────────────────────────
	http.HandleFunc("GET /api/v1/bookings/{id}/attendees",
		withAuth(h.ListAttendees))
	http.HandleFunc("POST /api/v1/bookings/{id}/attendees",
		withAuthAndRole(h.AddAttendee, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("PUT /api/v1/bookings/{id}/attendees/{attendee_id}",
		withAuthAndRole(h.UpdateAttendee, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
	http.HandleFunc("DELETE /api/v1/bookings/{id}/attendees/{attendee_id}",
		withAuthAndRole(h.RemoveAttendee, models.RoleBranchAdmin, models.RoleManager))
}

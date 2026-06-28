package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterCorProfileRoutes(
	corHandler *handlers.CorProfileHandler,
	bookingReqHandler *handlers.CorporateBookingRequestHandler,
) {
	// ─── Companies (backoffice) — admin only for mutations ────────────────────
	http.HandleFunc("GET /api/v1/clients/companies",
		withAuth(corHandler.ListCompanies))
	http.HandleFunc("GET /api/v1/clients/companies/{id}",
		withAuth(corHandler.GetCompany))
	http.HandleFunc("PUT /api/v1/clients/companies/{id}",
		withAuthAndRole(corHandler.UpdateCompany, models.RoleAdmin))

	// ─── Branches (backoffice) — admin only for mutations ─────────────────────
	http.HandleFunc("GET /api/v1/clients/companies/{id}/branches",
		withAuth(corHandler.ListBranches))
	http.HandleFunc("POST /api/v1/clients/companies/{id}/branches",
		withAuthAndRole(corHandler.CreateBranch, models.RoleAdmin))
	http.HandleFunc("PUT /api/v1/clients/companies/{id}/branches/{branch_id}",
		withAuthAndRole(corHandler.UpdateBranch, models.RoleAdmin))

	// ─── Profiles (backoffice) — admin only for mutations ─────────────────────
	http.HandleFunc("GET /api/v1/clients/companies/{id}/profiles",
		withAuth(corHandler.ListProfiles))
	http.HandleFunc("POST /api/v1/clients/companies/{id}/profiles",
		withAuthAndRole(corHandler.CreateProfile, models.RoleAdmin))
	http.HandleFunc("GET /api/v1/clients/profiles/{id}",
		withAuth(corHandler.GetProfile))
	http.HandleFunc("PUT /api/v1/clients/profiles/{id}",
		withAuthAndRole(corHandler.UpdateProfile, models.RoleAdmin))

	// ─── Corporate guests (backoffice) — admin only for mutations ────────────
	http.HandleFunc("GET /api/v1/clients/profiles/{id}/guests",
		withAuth(corHandler.ListGuests))
	http.HandleFunc("POST /api/v1/clients/profiles/{id}/guests",
		withAuthAndRole(corHandler.AddGuest, models.RoleAdmin))
	http.HandleFunc("PUT /api/v1/clients/profiles/{id}/guests/{guest_id}",
		withAuthAndRole(corHandler.UpdateGuest, models.RoleAdmin))
	http.HandleFunc("DELETE /api/v1/clients/profiles/{id}/guests/{guest_id}",
		withAuthAndRole(corHandler.DeleteGuest, models.RoleAdmin))

	// ─── Booking requests (web — authenticated submission) ───────────────────
	// Corporate accommodation submission: org_id and all company/approver data in request body
	http.HandleFunc("POST /api/v1/guest/bookings/corporate", withWebUserAuth(bookingReqHandler.SubmitAccommodation))
	http.HandleFunc("POST /api/v1/guest/bookings/corporate-event", withWebUserAuth(bookingReqHandler.SubmitEvent))
	http.HandleFunc("POST /api/v1/guest/bookings/corporate-meal", withWebUserAuth(bookingReqHandler.SubmitMeal))

	// ─── Booking requests (backoffice) ────────────────────────────────────────
	http.HandleFunc("GET /api/v1/booking-requests",
		withAuth(bookingReqHandler.List))
	http.HandleFunc("GET /api/v1/booking-requests/{id}",
		withAuth(bookingReqHandler.GetByID))
	http.HandleFunc("PUT /api/v1/booking-requests/{id}/approve",
		withAuthAndRole(bookingReqHandler.Approve,
			models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/booking-requests/{id}/reject",
		withAuthAndRole(bookingReqHandler.Reject,
			models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager))
	http.HandleFunc("PUT /api/v1/booking-requests/{id}/cancel",
		withAuthAndRole(bookingReqHandler.Cancel,
			models.RoleAdmin, models.RoleBranchAdmin, models.RoleManager, models.RoleReceptionist))
}

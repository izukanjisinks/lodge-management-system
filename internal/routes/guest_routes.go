package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterGuestRoutes(
	guestAuthHandler *handlers.GuestAuthHandler,
	guestBookingHandler *handlers.GuestBookingHandler,
) {
	// Public — self-registration (no auth required)
	http.HandleFunc("POST /api/v1/guest/register", withPublic(guestAuthHandler.Register))

	// Authenticated guest — profile
	http.HandleFunc("GET /api/v1/guest/me",
		withAuthAndRole(guestAuthHandler.Me, models.RoleGuest))

	http.HandleFunc("PUT /api/v1/guest/me",
		withAuthAndRole(guestAuthHandler.UpdateProfile, models.RoleGuest))

	// Authenticated guest — bookings
	http.HandleFunc("POST /api/v1/guest/bookings",
		withAuthAndRole(guestBookingHandler.Create, models.RoleGuest))

	http.HandleFunc("GET /api/v1/guest/bookings",
		withAuthAndRole(guestBookingHandler.List, models.RoleGuest))

	http.HandleFunc("GET /api/v1/guest/bookings/{id}",
		withAuthAndRole(guestBookingHandler.GetByID, models.RoleGuest))

	http.HandleFunc("PATCH /api/v1/guest/bookings/{id}/cancel",
		withAuthAndRole(guestBookingHandler.Cancel, models.RoleGuest))
}

package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
)

func RegisterGuestRoutes(
	guestAuthHandler *handlers.GuestAuthHandler,
	guestBookingHandler *handlers.GuestBookingHandler,
) {
	// Public — no auth required
	http.HandleFunc("POST /api/v1/guest/auth/register", withPublic(guestAuthHandler.Register))
	http.HandleFunc("POST /api/v1/guest/auth/login", withPublic(guestAuthHandler.Login))

	// Authenticated guest — profile
	http.HandleFunc("GET /api/v1/guest/me", withGuestAuth(guestAuthHandler.Me))
	http.HandleFunc("PUT /api/v1/guest/me", withGuestAuth(guestAuthHandler.UpdateProfile))

	// Authenticated guest — bookings
	http.HandleFunc("POST /api/v1/guest/bookings", withGuestAuth(guestBookingHandler.Create))
	http.HandleFunc("GET /api/v1/guest/bookings", withGuestAuth(guestBookingHandler.List))
	http.HandleFunc("GET /api/v1/guest/bookings/{id}", withGuestAuth(guestBookingHandler.GetByID))
	http.HandleFunc("PATCH /api/v1/guest/bookings/{id}/cancel", withGuestAuth(guestBookingHandler.Cancel))
}

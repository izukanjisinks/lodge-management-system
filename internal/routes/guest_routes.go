package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
)

func RegisterGuestRoutes(
	guestAuthHandler *handlers.GuestAuthHandler,
	roomHandler *handlers.RoomHandler,
	menuHandler *handlers.MenuHandler,
	venueHandler *handlers.VenueHandler,
) {
	// Public — no auth required
	http.HandleFunc("POST /api/v1/guest/auth/register", withPublic(guestAuthHandler.Register))
	http.HandleFunc("POST /api/v1/guest/auth/login", withPublic(guestAuthHandler.Login))
	http.HandleFunc("POST /api/v1/guest/auth/reset-password", withPublic(guestAuthHandler.ResetPassword))
	http.HandleFunc("GET /api/v1/guest/lodges", withPublic(guestAuthHandler.ListLodges))
	http.HandleFunc("GET /api/v1/guest/lodges/{org_id}", withPublic(guestAuthHandler.GetLodge))
	http.HandleFunc("GET /api/v1/guest/rooms", withPublic(roomHandler.GuestList))
	http.HandleFunc("GET /api/v1/guest/rooms/available", withPublic(roomHandler.GuestListAvailable))
	http.HandleFunc("GET /api/v1/guest/rooms/{id}", withPublic(roomHandler.GuestGetByID))
	http.HandleFunc("GET /api/v1/guest/menu", withPublic(menuHandler.GuestGetMenu))
	http.HandleFunc("GET /api/v1/guest/venues", withPublic(venueHandler.GuestList))
	http.HandleFunc("GET /api/v1/guest/venues/{id}", withPublic(venueHandler.GuestGetByID))

	// Authenticated guest — profile
	http.HandleFunc("GET /api/v1/guest/me", withGuestAuth(guestAuthHandler.Me))
	http.HandleFunc("PUT /api/v1/guest/me", withGuestAuth(guestAuthHandler.UpdateProfile))
}

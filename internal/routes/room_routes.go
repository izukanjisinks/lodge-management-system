package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterRoomRoutes(h *handlers.RoomHandler) {
	// Read — public (website guests browse rooms without logging in)
	http.HandleFunc("GET /api/v1/rooms",
		withPublic(h.List))

	http.HandleFunc("GET /api/v1/rooms/available",
		withPublic(h.ListAvailable))

	http.HandleFunc("GET /api/v1/rooms/{id}",
		withPublic(h.GetByID))

	// Write — admin and manager only
	http.HandleFunc("POST /api/v1/rooms",
		withAuthAndRole(h.Create, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/rooms/{id}",
		withAuthAndRole(h.Update, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("PATCH /api/v1/rooms/{id}/availability",
		withAuthAndRole(h.SetAvailability, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("DELETE /api/v1/rooms/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin, models.RoleManager))
}

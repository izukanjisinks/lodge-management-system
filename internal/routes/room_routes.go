package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterRoomRoutes(h *handlers.RoomHandler) {
	// Read — staff only (org-scoped via JWT)
	http.HandleFunc("GET /api/v1/rooms",
		withAuth(h.List))

	http.HandleFunc("GET /api/v1/rooms/available",
		withAuth(h.ListAvailable))

	http.HandleFunc("GET /api/v1/rooms/{id}",
		withAuth(h.GetByID))

	// Write — admin and manager only
	http.HandleFunc("POST /api/v1/rooms",
		withAuthAndRole(h.Create, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/rooms/{id}",
		withAuthAndRole(h.Update, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("PUT /api/v1/rooms/{id}/images",
		withAuthAndRole(h.UpdateImages, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("PATCH /api/v1/rooms/{id}/availability",
		withAuthAndRole(h.SetAvailability, models.RoleAdmin, models.RoleManager))

	http.HandleFunc("DELETE /api/v1/rooms/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin, models.RoleManager))
}

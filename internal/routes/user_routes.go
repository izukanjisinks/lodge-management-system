package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterUserRoutes(h *handlers.UserHandler) {
	http.HandleFunc("GET /api/v1/admin/users",
		withAuthAndRole(h.GetAll, models.RoleAdmin))

	http.HandleFunc("GET /api/v1/admin/users/{id}",
		withAuthAndRole(h.GetByID, models.RoleAdmin))

	http.HandleFunc("POST /api/v1/admin/users/{id}/role",
		withAuthAndRole(h.ChangeRole, models.RoleAdmin))

	http.HandleFunc("DELETE /api/v1/admin/users/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin))

	http.HandleFunc("POST /api/v1/admin/users/{id}/lock",
		withAuthAndRole(h.Lock, models.RoleAdmin))

	http.HandleFunc("POST /api/v1/admin/users/{id}/unlock",
		withAuthAndRole(h.Unlock, models.RoleAdmin))

	http.HandleFunc("GET /api/v1/profile",
		withAuth(h.GetProfile))
}

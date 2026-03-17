package routes

import (
	"net/http"

	"lodge-system/internal/handlers"
	"lodge-system/internal/models"
)

func RegisterUserRoutes(h *handlers.UserHandler) {
	// User management — frontend calls /api/v1/users (no /admin prefix)
	http.HandleFunc("POST /api/v1/users",
		withAuthAndRole(h.Create, models.RoleAdmin))

	http.HandleFunc("GET /api/v1/users",
		withAuthAndRole(h.GetAll, models.RoleAdmin))

	http.HandleFunc("GET /api/v1/users/{id}",
		withAuthAndRole(h.GetByID, models.RoleAdmin))

	http.HandleFunc("PUT /api/v1/users/{id}",
		withAuthAndRole(h.Update, models.RoleAdmin))

	http.HandleFunc("DELETE /api/v1/users/{id}",
		withAuthAndRole(h.Delete, models.RoleAdmin))

	http.HandleFunc("POST /api/v1/users/{id}/lock",
		withAuthAndRole(h.Lock, models.RoleAdmin))

	http.HandleFunc("POST /api/v1/users/{id}/unlock",
		withAuthAndRole(h.Unlock, models.RoleAdmin))

	// Keep /admin aliases for backward compatibility with Postman/other clients
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

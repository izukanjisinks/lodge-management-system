package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterUserRoutes(h *handlers.UserHandler) {
	// Get current user's profile - any authenticated user
	http.HandleFunc("GET /api/v1/users/profile",
		withAuth(h.GetProfile))

	// List all users - requires SuperAdmin or HRManager
	http.HandleFunc("GET /api/v1/users",
		withAuthAndRole(h.GetAll, models.RoleSuperAdmin, models.RoleManager, models.RoleHRManager))

	// Get user by ID - requires SuperAdmin or HRManager
	http.HandleFunc("GET /api/v1/users/{id}",
		withAuthAndRole(h.GetByID, models.RoleSuperAdmin, models.RoleManager, models.RoleHRManager))

	// Change user role - requires SuperAdmin
	http.HandleFunc("PATCH /api/v1/users/{id}/role",
		withAuthAndRole(h.ChangeRole, models.RoleSuperAdmin))

	// Delete user - requires SuperAdmin
	http.HandleFunc("DELETE /api/v1/users/{id}",
		withAuthAndRole(h.Delete, models.RoleSuperAdmin))

	// Lock user account - requires SuperAdmin or HRManager
	http.HandleFunc("POST /api/v1/users/{id}/lock",
		withAuthAndRole(h.Lock, models.RoleSuperAdmin, models.RoleHRManager))

	// Unlock user account - requires SuperAdmin or HRManager
	http.HandleFunc("POST /api/v1/users/{id}/unlock",
		withAuthAndRole(h.Unlock, models.RoleSuperAdmin, models.RoleHRManager))
}

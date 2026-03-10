package routes

import (
	"net/http"

	"hr-system/internal/handlers"
	"hr-system/internal/models"
)

func RegisterAuthRoutes(authHandler *handlers.AuthHandler) {
	http.HandleFunc("GET /api/v1/auth/me",
		withAuth(authHandler.Me))

	http.HandleFunc("POST /api/v1/auth/logout",
		withAuth(authHandler.Logout))

	http.HandleFunc("POST /api/v1/auth/register",
		withAuthAndRole(authHandler.Register, models.RoleSuperAdmin, models.RoleHRManager))
}
